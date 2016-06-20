package com.danielfireman.phd;

import java.lang.management.GarbageCollectorMXBean;
import java.lang.management.ManagementFactory;
import java.lang.management.OperatingSystemMXBean;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

import org.jooby.Jooby;
import org.jooby.Results;
import org.jooby.json.Jackson;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class App extends Jooby {
	{
		use(new Jackson());
	
		final AtomicInteger reqCount = new AtomicInteger(0);
		final AtomicInteger reqCompleted = new AtomicInteger(0);

		complete("/echo", (req, resp, cause) -> {
			reqCompleted.incrementAndGet();
		});

		before("/echo", (req, resp) -> {
			reqCount.incrementAndGet();
		});
		
		post("/echo", (req, resp) -> {
			resp.send(Results.json(req.body().to(Message.class)));
		});
				
		final Logger logger = LoggerFactory.getLogger("mon");
		final ScheduledExecutorService memSchedulerLogger = Executors.newScheduledThreadPool(1);
		memSchedulerLogger.scheduleAtFixedRate(new Runnable() {
			@Override
			public void run() {
				for (GarbageCollectorMXBean bean : ManagementFactory.getGarbageCollectorMXBeans()) {
					// colcount: total number of collections that have occurred.
					// This method returns -1 if the collection count is
					// undefined for this collector.

					// coltime: approximate accumulated collection elapsed time
					// in milliseconds. This method returns -1 if the collection
					// elapsed time is undefined for this collector.
					logger.info(String.format("gc name:%s colcount:%s coltime:%s", bean.getName(),
							bean.getCollectionCount(), bean.getCollectionTime()));
				}

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
				logger.info(String.format("cpu systemload:%s nproc:%s", osBean.getSystemLoadAverage(), osBean.getAvailableProcessors()));

				logger.info(String.format("req count:%s completed:%s", reqCount, reqCompleted));
			}
		}, 0, 30, TimeUnit.SECONDS);
	}
	
	public static void main(final String[] args) throws Throwable {
		run(App::new, args);
	}
}
