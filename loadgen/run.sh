#!/bin/bash
set -x
killall loadgen
rm logs/*

# This order must be in sync with run_multi.sh
CLIENT_ID=$1
NUM_ROUNDS=$2
SERVER_ADDR=$3

WARMUP_STEPS=3
TIMEOUT=100ms
STEP_DURATION=5s
INITIAL_QPS=10
STEP_SIZE=5
MAX_QPS=250

for i in `seq 1 ${NUM_ROUNDS}`
do
	sleep 1.5m
	./loadgen \
--num_warmup_steps=${WARMUP_STEPS} \
--initial_qps=${INITIAL_QPS} \
--step_size=${STEP_SIZE} \
--timeout=${TIMEOUT} \
--step_duration=${STEP_DURATION} \
--max_qps=${MAX_QPS} \
--addr=${SERVER_ADDR} > logs/client_${i}_${CLIENT_ID}
done
