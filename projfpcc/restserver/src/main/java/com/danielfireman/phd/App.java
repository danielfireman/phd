package com.danielfireman.phd;

import com.codahale.metrics.*;
import com.codahale.metrics.jvm.GarbageCollectorMetricSet;
import com.codahale.metrics.jvm.MemoryUsageGaugeSet;
import com.codahale.metrics.jvm.ThreadStatesGaugeSet;
import com.google.common.collect.Lists;
import com.sun.management.GarbageCollectionNotificationInfo;
import org.jooby.Jooby;
import org.jooby.Results;
import org.jooby.Status;
import org.jooby.metrics.Metrics;

import javax.management.Notification;
import javax.management.NotificationEmitter;
import javax.management.NotificationListener;
import javax.management.openmbean.CompositeData;
import java.io.File;
import java.lang.management.GarbageCollectorMXBean;
import java.lang.management.ManagementFactory;
import java.lang.management.MemoryPoolMXBean;
import java.lang.management.MemoryUsage;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicLong;

public class App extends Jooby {

    private static final int LOG_INTERVAL_SECS = 5;
    RequestCounter counter = new RequestCounter();

    {
        installGCMonitoring();
        String suffix = System.getenv("ROUND") != null ? "_" + System.getenv("ROUND") : "";
        use(new Metrics()
                .request()
                .metric("memory" + suffix, new MemoryUsageGaugeSet())
                .metric("threads" + suffix, new ThreadStatesGaugeSet())
                .metric("gc" + suffix, new GarbageCollectorMetricSet())
                .metric("fgc" + suffix, new ForcedGCMetricSet(counter))
                .metric("cpu" + suffix, new CpuInfoGaugeSet())
                .reporter(registry -> {
                    CsvReporter reporter = CsvReporter.forRegistry(registry)
                            .convertRatesTo(TimeUnit.SECONDS)
                            .convertDurationsTo(TimeUnit.MILLISECONDS)
                            .build(new File("logs/"));
                    reporter.start(LOG_INTERVAL_SECS, TimeUnit.SECONDS);
                    return reporter;
                }));

        if (System.getenv("CONTROL_GC") != null && !System.getenv("CONTROL_GC").equals("")) {
            double threshold = Double.parseDouble(System.getenv("CONTROL_GC"));
            System.out.println("Controlling GC wih threshold: " + threshold);
            ExecutorService gcPool = Executors.newSingleThreadExecutor();

            use("GET", "*", (req, rsp, chain) -> {
                if (counter.doingGC.get()) {
                    rsp.header("Retry-After", getRetryAfter(counter.unavailabilityHist.getSnapshot(), counter.unavailabilityStartTime.get())).status(Status.TOO_MANY_REQUESTS)
                            .length(0)
                            .end();
                    return;
                }
                if (counter.incoming.get() % counter.sampleRate.get() == 0) {
                    synchronized (counter) {
                        if (counter.doingGC.get()) {
                            rsp.header("Retry-After", getRetryAfter(counter.unavailabilityHist.getSnapshot(), counter.unavailabilityStartTime.get())).status(Status.TOO_MANY_REQUESTS)
                                    .length(0)
                                    .end();
                            return;
                        }
                        MemoryUsage youngUsage = counter.youngPool.getUsage();
                        MemoryUsage oldUsage = counter.oldPool.getUsage();
                        if ((double) youngUsage.getUsed() / (double) youngUsage.getCommitted() > threshold ||
                                (double) oldUsage.getUsed() / (double) oldUsage.getCommitted() > threshold) {
                            counter.doingGC.set(true);
                            counter.unavailabilityStartTime.set(System.currentTimeMillis());
                            rsp.header("Retry-After", getRetryAfter(counter.unavailabilityHist.getSnapshot(), counter.unavailabilityStartTime.get())).status(Status.TOO_MANY_REQUESTS)
                                    .length(0)
                                    .end();

                            // This should return immediately.
                            gcPool.execute(()-> {
                                // Calculating next sample rate.
                                // The main idea is to get 1/10th of the requests that arrived since last GC and bound
                                // this number to [10, 300].
                                Snapshot sRH = counter.requestTimeHistogram.getSnapshot();
                                System.out.println("ReqHist: " + sRH.getMedian() + " " + sRH.get95thPercentile() + " " + sRH.get99thPercentile());
                                counter.sampleRate.set(Math.min(500, Math.max(50L, (long) ((double)counter.incoming.get()/5d))));

                                // Loop waiting for the queue to get empty. Each iteration waits the median of request
                                // processing time.
                                long waitTime = (long) counter.requestTimeHistogram.getSnapshot().getMedian();
                                while (counter.finished.get() < counter.incoming.get()) {
                                    try {
                                        Thread.sleep(waitTime);
                                    } catch (InterruptedException ie) {
                                        throw new RuntimeException(ie);
                                    }
                                }

                                // Finally GC.
                                counter.gcCountForcedGC.incrementAndGet();
                                long startTime = System.currentTimeMillis();
                                System.gc();
                                counter.gcTimeForcedGCMillis.addAndGet(System.currentTimeMillis() - startTime);

                                // Now we can start doing GC and attend new requests.
                                counter.unavailabilityHist.update(System.currentTimeMillis() - counter.unavailabilityStartTime.get());
                                counter.incoming.set(0);
                                counter.finished.set(0);
                                counter.doingGC.set(false);
                            });
                            return;
                        }
                    }
                }
                long startTime = System.currentTimeMillis();
                counter.incoming.incrementAndGet();
                chain.next(req, rsp);
                counter.finished.incrementAndGet();
                counter.requestTimeHistogram.update(System.currentTimeMillis() - startTime);
            });
        }

        get("/quit", () -> {
            System.exit(0);
            return Results.ok();
        });

        get("/numprimes/:max", (req) -> {
            long startTime = System.currentTimeMillis();
            long max = req.param("max").longValue();
            long count = 0;
            for (long i = 3; i <= max; i++) {
                boolean isPrime = true;
                for (long j = 2; j <= i / 2 && isPrime; j++) {
                    isPrime = i % j > 0;
                }
                if (isPrime) {
                    count++;
                }
            }
            long elapsed = System.currentTimeMillis() - startTime;
            return Results.ok(count + "," + elapsed);
        });

        get("/allocmem/:amount", (req) -> {
            int arraySize = req.param("amount").intValue();
            byte[] array = new byte[arraySize];
            return Results.ok();
        });

        get("/allocandhold/:amount/:millis", (req) -> {
            int arraySize = req.param("amount").intValue();
            int millis = req.param("millis").intValue();
            byte[] array = new byte[arraySize];
            Thread.sleep(millis);
            return Results.ok();
        });
    }

    public static void main(final String[] args) throws Throwable {
        run(App::new, args);
    }

    public static void installGCMonitoring() {
        //get all the GarbageCollectorMXBeans - there's one for each heap generation
        //so probably two - the old generation and young generation
        List<GarbageCollectorMXBean> gcbeans = java.lang.management.ManagementFactory.getGarbageCollectorMXBeans();
        //Install a notifcation handler for each bean
        for (GarbageCollectorMXBean gcbean : gcbeans) {
            NotificationEmitter emitter = (NotificationEmitter) gcbean;
            //use an anonymously generated listener for this example
            // - proper code should really use a named class
            NotificationListener listener = new NotificationListener() {
                //keep a count of the total time spent in GCs
                long totalGcDuration = 0;

                //implement the notifier callback handler
                @Override
                public void handleNotification(Notification notification, Object handback) {
                    //we only handle GARBAGE_COLLECTION_NOTIFICATION notifications here
                    if (notification.getType().equals(GarbageCollectionNotificationInfo.GARBAGE_COLLECTION_NOTIFICATION)) {
                        //get the information associated with this notification
                        GarbageCollectionNotificationInfo info = GarbageCollectionNotificationInfo.from((CompositeData) notification.getUserData());
                        //get all the info and pretty print it
                        long duration = info.getGcInfo().getDuration();
                        String gctype = info.getGcAction();
                        if ("end of minor GC".equals(gctype)) {
                            gctype = "Young Gen GC";
                        } else if ("end of major GC".equals(gctype)) {
                            gctype = "Old Gen GC";
                        }
                        System.out.println();
                        System.out.println(gctype + ": - " + info.getGcInfo().getId() + " " + info.getGcName() + " (from " + info.getGcCause() + ") " + duration + " milliseconds; start-end times " + info.getGcInfo().getStartTime() + "-" + info.getGcInfo().getEndTime());
                        //System.out.println("GcInfo CompositeType: " + info.getGcInfo().getCompositeType());
                        //System.out.println("GcInfo MemoryUsageAfterGc: " + info.getGcInfo().getMemoryUsageAfterGc());
                        //System.out.println("GcInfo MemoryUsageBeforeGc: " + info.getGcInfo().getMemoryUsageBeforeGc());

                        //Get the information about each memory space, and pretty print it
                        Map<String, MemoryUsage> membefore = info.getGcInfo().getMemoryUsageBeforeGc();
                        Map<String, MemoryUsage> mem = info.getGcInfo().getMemoryUsageAfterGc();
                        for (Map.Entry<String, MemoryUsage> entry : mem.entrySet()) {
                            String name = entry.getKey();
                            MemoryUsage memdetail = entry.getValue();
                            long memInit = memdetail.getInit();
                            long memCommitted = memdetail.getCommitted();
                            long memMax = memdetail.getMax();
                            long memUsed = memdetail.getUsed();
                            MemoryUsage before = membefore.get(name);
                            long beforepercent = ((before.getUsed() * 1000L) / before.getCommitted());
                            long percent = ((memUsed * 1000L) / before.getCommitted()); //>100% when it gets expanded

                            System.out.print(name + (memCommitted == memMax ? "(fully expanded)" : "(still expandable)") + "used: " + (beforepercent / 10) + "." + (beforepercent % 10) + "%->" + (percent / 10) + "." + (percent % 10) + "%(" + ((memUsed / 1048576) + 1) + "MB) / ");
                        }
                        System.out.println();
                        totalGcDuration += info.getGcInfo().getDuration();
                        long percent = totalGcDuration * 1000L / info.getGcInfo().getEndTime();
                        System.out.println("GC cumulated overhead " + (percent / 10) + "." + (percent % 10) + "%");
                    }
                }
            };

            //Add the listener
            emitter.addNotificationListener(listener, null, null);
        }
    }

    static class RequestCounter {
        AtomicLong incoming = new AtomicLong();
        AtomicLong finished = new AtomicLong();
        AtomicLong sampleRate = new AtomicLong(10);
        MemoryPoolMXBean youngPool;
        MemoryPoolMXBean oldPool;
        AtomicBoolean doingGC = new AtomicBoolean(false);
        AtomicLong gcCountForcedGC = new AtomicLong(0);
        AtomicLong gcTimeForcedGCMillis = new AtomicLong(0);
        AtomicLong unavailabilityStartTime = new AtomicLong(0);
        Histogram unavailabilityHist = new Histogram(new SlidingWindowReservoir(10));
        Histogram requestTimeHistogram = new Histogram(new SlidingWindowReservoir(300));

        RequestCounter() {
            for (final MemoryPoolMXBean pool : ManagementFactory.getMemoryPoolMXBeans()) {
                if (pool.getName().contains("Eden")) {
                    youngPool = pool;
                    continue;
                }
                if (pool.getName().contains("Old")) {
                    oldPool = pool;
                    continue;
                }
            }
        }
    }

    static String getRetryAfter(Snapshot s, long lastGCStartTime) {
        long delta = System.currentTimeMillis() - lastGCStartTime;
        return Double.toString((double) Math.max(0, (s.getMedian() + s.getStdDev() - delta) / 1000d));
    }

    static class ForcedGCMetricSet implements MetricSet {
        private final RequestCounter counter;

        ForcedGCMetricSet(RequestCounter counter) {
            this.counter = counter;
        }

        @Override
        public Map<String, Metric> getMetrics() {
            final Map<String, Metric> gauges = new HashMap<>();
            gauges.put("count", (Gauge<Long>) () -> counter.gcCountForcedGC.get());
            gauges.put("time", (Gauge<Long>) () -> counter.gcTimeForcedGCMillis.get());
            return gauges;
        }
    }
}
