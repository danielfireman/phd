#!/bin/bash

set -x

rm logs/*

# This order must be in sync with run_multi.sh
CLIENT_ID=$1
NUM_ROUNDS=$2
WARMUP_STEPS=$3
TIMEOUT=$4
STEP_DURATION=$5
INITIAL_QPS=$6
STEP_SIZE=$7
MAX_QPS=$8
SUFFIXES=$9
KEEP_DURATION=${10}
SERVER=${11}

for i in `seq 1 ${NUM_ROUNDS}`
do
    if [ "${NUM_ROUNDS}" -gt "1" ]; then
	    sleep 60s
	fi
	./loadgen \
--num_warmup_steps=${WARMUP_STEPS} \
--initial_qps=${INITIAL_QPS} \
--step_size=${STEP_SIZE} \
--timeout=${TIMEOUT} \
--step_duration=${STEP_DURATION} \
--max_qps=${MAX_QPS} \
--msg_suffixes=${SUFFIXES} \
--keep_duration=${KEEP_DURATION} \
--addr=${SERVER} > logs/client_${i}_${CLIENT_ID}
done
