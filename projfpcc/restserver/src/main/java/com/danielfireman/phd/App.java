package com.danielfireman.phd;

import java.lang.management.GarbageCollectorMXBean;
import java.lang.management.ManagementFactory;
import java.lang.management.OperatingSystemMXBean;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

import org.jooby.Jooby;
import org.jooby.Results;
import org.jooby.json.Jackson;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class App extends Jooby {
	{
		use(new Jackson());

		onStart(()->startLogger());

		get("/msg", () -> {
			Message msg = new Message();
			msg.content = "My super big top ultra big content";
			return Results.json(msg);
		});
	}

	private static void startLogger() {
		final Logger gcLogger = LoggerFactory.getLogger("gc");
		gcLogger.info("timestamp,name,colcount,coltime");

		final Logger cpuLogger = LoggerFactory.getLogger("cpu");
		cpuLogger.info("timestamp,systemload,nprocs");

		final ScheduledExecutorService memSchedulerLogger = Executors.newScheduledThreadPool(1);
		memSchedulerLogger.scheduleAtFixedRate(new Runnable() {
			@Override
			public void run() {
				long time = System.currentTimeMillis();

				// ## GC ##
				for (GarbageCollectorMXBean bean : ManagementFactory.getGarbageCollectorMXBeans()) {
					// colcount: total number of collections that have occurred.
					// This method returns -1 if the collection count is
					// undefined for this collector.

					// coltime: approximate accumulated collection elapsed time
					// in milliseconds. This method returns -1 if the collection
					// elapsed time is undefined for this collector.
					gcLogger.info(String.format("%s,\"%s\",%s,%s", time, bean.getName(), bean.getCollectionCount(),
							bean.getCollectionTime()));
				}

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
				OperatingSystemMXBean osBean = ManagementFactory.getPlatformMXBean(OperatingSystemMXBean.class);
				cpuLogger.info(String.format("%s,%s,%s", time, osBean.getSystemLoadAverage(), osBean.getAvailableProcessors()));
			}
		}, 0, 10, TimeUnit.SECONDS);
	}

	public static void main(final String[] args) throws Throwable {
		run(App::new, args);
	}
}
