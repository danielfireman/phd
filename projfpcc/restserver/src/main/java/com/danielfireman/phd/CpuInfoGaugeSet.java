package com.danielfireman.phd;

import com.codahale.metrics.Gauge;
import com.codahale.metrics.Metric;
import com.codahale.metrics.MetricSet;
import com.google.common.collect.ImmutableMap;
import com.sun.management.OperatingSystemMXBean;

import java.lang.management.ManagementFactory;
import java.util.Map;

@SuppressWarnings("restriction")
public class CpuInfoGaugeSet implements MetricSet {
    private OperatingSystemMXBean bean;

    public CpuInfoGaugeSet() {
        bean = (OperatingSystemMXBean) ManagementFactory.getOperatingSystemMXBean();
    }

    @Override
    public Map<String, Metric> getMetrics() {
        return ImmutableMap.<String, Metric>builder()
                .put("count", (Gauge<Integer>) () -> bean.getAvailableProcessors())
                .put("load", (Gauge<Double>) () -> bean.getProcessCpuLoad())
                .put("time", (Gauge<Long>) () -> bean.getProcessCpuTime())
                .build();
    }
}
