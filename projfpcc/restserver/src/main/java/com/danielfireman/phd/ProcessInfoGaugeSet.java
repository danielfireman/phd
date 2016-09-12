package com.danielfireman.phd;

import com.codahale.metrics.Gauge;
import com.codahale.metrics.Metric;
import com.codahale.metrics.MetricSet;
import com.google.common.collect.ImmutableMap;
import oshi.SystemInfo;
import oshi.software.os.OSProcess;
import oshi.software.os.OperatingSystem;

import java.util.Map;

public class ProcessInfoGaugeSet implements MetricSet {

    private final OSProcess process;

    ProcessInfoGaugeSet() {
        OperatingSystem os = new SystemInfo().getOperatingSystem();
        process = os.getProcess(os.getProcessId());
    }

    @Override
    public Map<String, Metric> getMetrics() {
        return ImmutableMap.<String, Metric>builder()
                .put("kernel.time", (Gauge<Long>) () -> process.getKernelTime())
                .put("user.time", (Gauge<Long>) () -> process.getUserTime())
                .put("up.time", (Gauge<Long>) () -> process.getUpTime())
                .build();
    }
}
