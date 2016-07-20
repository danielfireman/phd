#!/bin/bash
set -x

declare -a CLIENTS=('10.4.5.132' '10.4.5.130' '10.4.5.133' '10.4.5.134')

NUM_ROUNDS="$1"
SERVER="http://10.4.2.103:8080"
NUM_CLIENTS=2
ACTIVE_CLIENTS=0

for CLIENT in ${CLIENTS[@]};
do
	if [ $ACTIVE_CLIENTS == $NUM_CLIENTS ]; then
		break
	fi
	ssh -i ~/fireman.sururu.key ubuntu@${CLIENT} "cd ~/phd/loadgen;./run.sh ${NUM_CLIENTS} ${NUM_ROUNDS} ${SERVER}" &
	ACTIVE_CLIENTS=`expr $ACTIVE_CLIENTS + 1`
done
