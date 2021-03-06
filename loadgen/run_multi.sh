#!/bin/bash
set -x

source configrc

killall loadgen

NUM_ROUNDS="$1"
ACTIVE_CLIENT=1

for CLIENT in ${CLIENTS[@]};
do
	port=$(getport ${CLIENT})
	ssh -i ~/fireman.sururu.key ${SSH_ADDR} -p ${port} "cd ~/phd/loadgen;./run.sh \
	${ACTIVE_CLIENT} \
	${NUM_ROUNDS} \
	${WARMUP_STEPS} \
    ${TIMEOUT} \
    ${STEP_DURATION} \
    ${INITIAL_QPS} \
    ${STEP_SIZE} \
    ${MAX_QPS} \
    ${SUFFIXES} \
    ${KEEP_DURATION} \
    ${SERVER}" &
	if [ ${ACTIVE_CLIENT} == ${NUM_CLIENTS} ]; then
		break
	fi
	ACTIVE_CLIENT=`expr ${ACTIVE_CLIENT} + 1`
done
