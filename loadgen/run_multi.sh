#!/bin/bash
set -x

declare -a CLIENTS=('10.4.5.132' '10.4.5.130' '10.4.5.133' '10.4.5.134')

NUM_ROUNDS="$1"
SERVER="http://10.4.5.216:8080"
NUM_CLIENTS=2
ACTIVE_CLIENT=1

for CLIENT in ${CLIENTS[@]};
do
	ssh -i ~/fireman.sururu.key ubuntu@${CLIENT} "cd ~/phd/loadgen;./run.sh ${ACTIVE_CLIENT} ${NUM_ROUNDS} ${SERVER}" &
	if [ $ACTIVE_CLIENT == $NUM_CLIENTS ]; then
		break
	fi
	ACTIVE_CLIENT=`expr $ACTIVE_CLIENTS + 1`
done
