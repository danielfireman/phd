package com.danielfireman.phd;

import java.io.File;
import java.lang.management.GarbageCollectorMXBean;
import java.lang.management.ManagementFactory;
import java.lang.management.MemoryPoolMXBean;
import java.lang.management.MemoryUsage;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicLong;

import com.sun.management.GarbageCollectionNotificationInfo;
import org.jooby.Jooby;
import org.jooby.Results;
import org.jooby.Status;
import org.jooby.metrics.Metrics;

import com.codahale.metrics.CsvReporter;
import com.codahale.metrics.jvm.GarbageCollectorMetricSet;
import com.codahale.metrics.jvm.MemoryUsageGaugeSet;
import com.codahale.metrics.jvm.ThreadStatesGaugeSet;

import javax.management.Notification;
import javax.management.NotificationEmitter;
import javax.management.NotificationListener;
import javax.management.openmbean.CompositeData;

public class App extends Jooby {
    static class RequestCounter {
        AtomicLong incoming = new AtomicLong();
        AtomicLong finished = new AtomicLong();
        AtomicLong numReqAtLastGC = new AtomicLong();
        AtomicLong sampleRate = new AtomicLong(10);
        List<MemoryPoolMXBean> mem = ManagementFactory.getMemoryPoolMXBeans();
        AtomicBoolean doingGC = new AtomicBoolean(false);
    }

    static RequestCounter counter = new RequestCounter();

	private static final int LOG_INTERVAL_SECS = 5;
	{
        installGCMonitoring();
		String suffix = System.getenv("ROUND") != null ? "_" + System.getenv("ROUND") : "";
		use(new Metrics().request().threadDump().metric("memory" + suffix, new MemoryUsageGaugeSet())
				.request()
				.metric("threads" + suffix, new ThreadStatesGaugeSet())
				.metric("gc" + suffix, new GarbageCollectorMetricSet())
				.metric("cpu" + suffix, new CpuInfoGaugeSet())
				.reporter(registry -> {
					CsvReporter reporter = CsvReporter.forRegistry(registry)
                            .convertRatesTo(TimeUnit.SECONDS)
							.convertDurationsTo(TimeUnit.MILLISECONDS)
                            .build(new File("logs/"));
					reporter.start(LOG_INTERVAL_SECS, TimeUnit.SECONDS);
					return reporter;
				}));

        if (System.getenv("CONTROL_GC") != null) {
            use("GET", "/numprimes/:max", (req, rsp, chain) -> {
                if (counter.doingGC.get()) {
                            rsp.status(Status.TOO_MANY_REQUESTS)
			        .length(0)
			        .end();

                    return;
                }

                boolean doGC = false;
                String cause = "";
                if (counter.incoming.get() % counter.sampleRate.get() == 0) {
                    synchronized (counter) {
                        // double checked locking
                        if (counter.doingGC.get()) {
                            rsp.status(Status.TOO_MANY_REQUESTS)
			        .length(0)
			        .end();
                            return;
                        }
                        for (final MemoryPoolMXBean pool : counter.mem) {
                            double perc = (double) pool.getUsage().getUsed() / (double) pool.getUsage().getCommitted();
                            String name = pool.getName();
                            if ((name.contains("Eden") || name.contains("Old")) && perc > 0.75) {
                                cause = name;
                                counter.doingGC.set(true);
                                doGC = true;
                                break;
                            }
                        }
                    }
                }

                if (doGC) {
                    long inc = counter.incoming.get();
                    long numReqLast = counter.numReqAtLastGC.get();
                    counter.sampleRate.set(Math.min(300, Math.max(10L, (long) ((double) (inc - numReqLast) / 10d))));
                    counter.numReqAtLastGC.set(inc);
                            rsp.status(Status.TOO_MANY_REQUESTS)
			        .length(0)
			        .end();

                    System.out.println("\n\nCause:" + cause + " | Incoming: " + counter.incoming + " Finished:" + counter.finished + " SampleRate: " + counter.sampleRate.get());
		    // Waiting until queue gets empty.
                    while (counter.finished.get() < counter.incoming.get()) {
                        Thread.currentThread().sleep(50);
                    }
                    System.gc();
                    counter.doingGC.set(false);
                } else {
                    counter.incoming.incrementAndGet();
                    chain.next(req, rsp);
                    counter.finished.incrementAndGet();
                }
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

	public static void installGCMonitoring(){
		//get all the GarbageCollectorMXBeans - there's one for each heap generation
		//so probably two - the old generation and young generation
		List<GarbageCollectorMXBean> gcbeans = java.lang.management.ManagementFactory.getGarbageCollectorMXBeans();
		//Install a notifcation handler for each bean
		for (GarbageCollectorMXBean gcbean : gcbeans) {
			System.out.println(gcbean);
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
						System.out.println(gctype + ": - " + info.getGcInfo().getId()+ " " + info.getGcName() + " (from " + info.getGcCause()+") "+duration + " milliseconds; start-end times " + info.getGcInfo().getStartTime()+ "-" + info.getGcInfo().getEndTime());
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
							long beforepercent = ((before.getUsed()*1000L)/before.getCommitted());
							long percent = ((memUsed*1000L)/before.getCommitted()); //>100% when it gets expanded

							System.out.print(name + (memCommitted==memMax?"(fully expanded)":"(still expandable)") +"used: "+(beforepercent/10)+"."+(beforepercent%10)+"%->"+(percent/10)+"."+(percent%10)+"%("+((memUsed/1048576)+1)+"MB) / ");
						}
						System.out.println();
						totalGcDuration += info.getGcInfo().getDuration();
						long percent = totalGcDuration*1000L/info.getGcInfo().getEndTime();
						System.out.println("GC cumulated overhead "+(percent/10)+"."+(percent%10)+"%");
					}
				}
			};

			//Add the listener
			emitter.addNotificationListener(listener, null, null);
		}
	}
}
