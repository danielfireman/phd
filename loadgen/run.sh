#!/bin/bash
set -x
killall main
rm logs/*

WARMUP_STEPS=3
ADDR=http://10.4.2.103:8080
TIMEOUT=30ms
STEP_DURATION=10s
INITIAL_QPS=50
STEP_SIZE=50
MAX_QPS=1500
for i in `seq 1 $1`
do
	sleep 1.5m
	go run main.go --num_warmup_steps=${WARMUP_STEPS} --initial_qps=${INITIAL_QPS} --step_size=${STEP_SIZE} --timeout=${TIMEOUT} --max_qps=${MAX_QPS} --addr=${ADDR} > logs/client_${i}_1 &
	go run main.go --num_warmup_steps=${WARMUP_STEPS} --initial_qps=${INITIAL_QPS} --step_size=${STEP_SIZE} --timeout=${TIMEOUT} --max_qps=${MAX_QPS} --addr=${ADDR} > logs/client_${i}_2
done
