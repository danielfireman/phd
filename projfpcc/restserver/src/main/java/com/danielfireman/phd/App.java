package com.danielfireman.phd;

import java.io.File;
import java.util.concurrent.TimeUnit;

import org.jooby.Jooby;
import org.jooby.Results;
import org.jooby.metrics.Metrics;

import com.codahale.metrics.CsvReporter;
import com.codahale.metrics.jvm.FileDescriptorRatioGauge;
import com.codahale.metrics.jvm.GarbageCollectorMetricSet;
import com.codahale.metrics.jvm.MemoryUsageGaugeSet;
import com.codahale.metrics.jvm.ThreadStatesGaugeSet;

public class App extends Jooby {
	private static final int LOG_INTERVA_SECS = 10;
	{
		String suffix = System.getenv("ROUND") != null ? "_" + System.getenv("ROUND") : "";
		use(new Metrics().request().threadDump().metric("memory" + suffix, new MemoryUsageGaugeSet())
				.metric("threads" + suffix, new ThreadStatesGaugeSet())
				.metric("gc" + suffix, new GarbageCollectorMetricSet())
				.metric("fs" + suffix, new FileDescriptorRatioGauge()).metric("cpu" + suffix, new CpuInfoGaugeSet())
				.reporter(registry -> {
					CsvReporter reporter = CsvReporter.forRegistry(registry).convertRatesTo(TimeUnit.SECONDS)
							.convertDurationsTo(TimeUnit.MILLISECONDS).build(new File("logs/"));
					reporter.start(LOG_INTERVA_SECS, TimeUnit.SECONDS);
					return reporter;
				}));

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

		get("allocmem/:amount", (req) -> {
            int arraySize = req.param("amount").intValue();
            byte[] array = new byte[arraySize];
            return Results.ok();
		});
	}

	public static void main(final String[] args) throws Throwable {
		run(App::new, args);
	}
}
