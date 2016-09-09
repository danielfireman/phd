#!/bin/bash
set -x

source configrc

NUM_ROUNDS="$1"
ACTIVE_CLIENT=1

for CLIENT in ${CLIENTS[@]};
do
	ssh -i ~/fireman.sururu.key ubuntu@${CLIENT} "cd ~/phd/loadgen;./run.sh ${ACTIVE_CLIENT} ${NUM_ROUNDS}" &
	if [ $ACTIVE_CLIENT == $NUM_CLIENTS ]; then
		break
	fi
	ACTIVE_CLIENT=`expr $ACTIVE_CLIENT + 1`
done
