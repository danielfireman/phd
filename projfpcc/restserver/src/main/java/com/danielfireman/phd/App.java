package com.danielfireman.phd;

import java.io.File;
import java.lang.management.ManagementFactory;
import java.util.Locale;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

import org.jooby.Jooby;
import org.jooby.Results;
import org.jooby.json.Jackson;
import org.jooby.metrics.Metrics;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.codahale.metrics.CsvReporter;
import com.codahale.metrics.jvm.FileDescriptorRatioGauge;
import com.codahale.metrics.jvm.GarbageCollectorMetricSet;
import com.codahale.metrics.jvm.MemoryUsageGaugeSet;
import com.codahale.metrics.jvm.ThreadStatesGaugeSet;
import com.sun.management.OperatingSystemMXBean;

@SuppressWarnings("restriction")
public class App extends Jooby {
	{
		use(new Jackson());

		onStart(() -> {
			startLogger();
		});

		use(new Metrics()
				.request()
				.threadDump()
				.metric("memory", new MemoryUsageGaugeSet())
				.metric("threads", new ThreadStatesGaugeSet())
				.metric("gc", new GarbageCollectorMetricSet())
				.metric("fs", new FileDescriptorRatioGauge())
				.reporter(registry -> {
					CsvReporter reporter = CsvReporter.forRegistry(registry)
							.convertRatesTo(TimeUnit.SECONDS)
                            .convertDurationsTo(TimeUnit.MILLISECONDS)
                            .build(new File("logs/"));
					reporter.start(30, TimeUnit.SECONDS);
				    return reporter;
				}));

		get("/quit", () -> {
			System.exit(0);
			return Results.ok();
		});

		get("/msg", () -> {
			Message msg = new Message();
			msg.content = "My super big top ultra big content";
			return Results.json(msg);
		});
	}

	private static void startLogger() {
		final Logger cpuLogger = LoggerFactory.getLogger("cpu");
		cpuLogger.info("timestamp,systemload,nprocs");

		final ScheduledExecutorService memSchedulerLogger = Executors.newScheduledThreadPool(1);
		memSchedulerLogger.scheduleAtFixedRate(new Runnable() {
			@Override
			public void run() {
				// ## CPU ##
				// systemload: the system load average for the last minute. The
				// system load average is the sum of the number of runnable
				// entities queued to the available processors and the number of
				// runnable entities running on the available processors
				// averaged over a period of time. The way in which the load
				// average is calculated is operating system specific but is
				// typically a damped time-dependent average.
				// Other reference:
				// http://blog.scoutapp.com/articles/2009/07/31/understanding-load-averages
				OperatingSystemMXBean osBean = (OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean();
				cpuLogger.info(String.format(Locale.US, "%s,%.2f,%s",
						System.currentTimeMillis(),
						osBean.getProcessCpuLoad(),
						osBean.getAvailableProcessors()));
			}
		}, 0, 10, TimeUnit.SECONDS);
	}

	public static void main(final String[] args) throws Throwable {
		run(App::new, args);
	}
}
