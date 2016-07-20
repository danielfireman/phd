#!/bin/bash
set -x
killall loadgen
rm logs/*

# This order must be in sync with run_multi.sh
NUM_CLIENTS=$1
NUM_ROUNDS=$2
SERVER_ADDR=$3

WARMUP_STEPS=3
TIMEOUT=30ms
STEP_DURATION=10s
INITIAL_QPS=50
STEP_SIZE=50
MAX_QPS=1200
for i in `seq 1 ${NUM_ROUNDS}`
do
#	sleep 1.5m
	id="`cat /var/lib/dbus/machine-id`"
	./loadgen \
--num_warmup_steps=${WARMUP_STEPS} \
--initial_qps=${INITIAL_QPS} \
--step_size=${STEP_SIZE} \
--timeout=${TIMEOUT} \
--max_qps=${MAX_QPS} \
--addr=${SERVER_ADDR} > logs/client_${NUM_CLIENTS}_${id}
done
