#!/bin/bash
set -x

if [[ -z "$1" ]]; then
	echo "Invalid execution id"
	exit -1
fi

source configrc

EXECUTION_ID=$1
ACTIVE_CLIENT=1

for CLIENT in ${CLIENTS[@]};
do
	fname="c${NUM_CORES}_${EXECUTION_ID}_${ACTIVE_CLIENT}.zip"
	ssh -i ~/fireman.sururu.key ubuntu@${CLIENT} "cd ~/phd/loadgen/logs; zip ${fname} *"
	scp -i ~/fireman.sururu.key ubuntu@${CLIENT}:~/phd/loadgen/logs/${fname} ~/expresults_fpcc/
	scp ~/expresults_fpcc/${fname} danielfireman@siri.lsd.ufcg.edu.br:/tmp/fireman_${fname}
	if [ $ACTIVE_CLIENT == $NUM_CLIENTS ]; then
		break
	fi
	ACTIVE_CLIENT=`expr $ACTIVE_CLIENT + 1`
done
