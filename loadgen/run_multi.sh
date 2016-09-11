#!/bin/bash
set -x

source configrc

NUM_ROUNDS="$1"
ACTIVE_CLIENT=1

for CLIENT in ${CLIENTS[@]};
do
	port=$(getport ${CLIENT})
	ssh -i ~/fireman.sururu.key ${SSH_ADDR} -p ${port} cd ~/phd/loadgen;./run.sh ${ACTIVE_CLIENT} ${NUM_ROUNDS} &
	if [ $ACTIVE_CLIENT == $NUM_CLIENTS ]; then
		break
	fi
	ACTIVE_CLIENT=`expr $ACTIVE_CLIENT + 1`
done
