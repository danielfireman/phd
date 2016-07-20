#!/bin/bash
set -x
killall main
rm logs/*

WARMUP_STEPS=3
TIMEOUT=30ms
STEP_DURATION=10s
INITIAL_QPS=50
STEP_SIZE=50
MAX_QPS=1000
for i in `seq 1 $1`
do
#	sleep 1.5m
	id="`cat /var/lib/dbus/machine-id`"
	./loadgen \
--num_warmup_steps=${WARMUP_STEPS} \
--initial_qps=${INITIAL_QPS} \
--step_size=${STEP_SIZE} \
--timeout=${TIMEOUT} \
--max_qps=${MAX_QPS} \
--addr=$2 > logs/client_${1}_{id}
done
