#!/bin/bash
set -x

if [[ -z "$1" ]]; then
	echo "Invalid execution id"
	exit -1
fi

source configrc

EXECUTION_ID=$1
ACTIVE_CLIENT=1

echo "Packing client logs"
for CLIENT in ${CLIENTS[@]};
do
	fname="c${NUM_CORES}_${EXECUTION_ID}_${ACTIVE_CLIENT}.zip"
	ssh -i ~/fireman.sururu.key ubuntu@${CLIENT} "cd ~/phd/loadgen/logs; zip ${fname} client*"
	scp -i ~/fireman.sururu.key ubuntu@${CLIENT}:~/phd/loadgen/logs/${fname} ~/expresults_fpcc/
	scp ~/expresults_fpcc/${fname} danielfireman@siri.lsd.ufcg.edu.br:/tmp/fireman_${fname}
	if [ $ACTIVE_CLIENT == $NUM_CLIENTS ]; then
		break
	fi
	ACTIVE_CLIENT=`expr $ACTIVE_CLIENT + 1`
done

echo "Packing server logs"
sfname="s${NUM_CORES}_${EXECUTION_ID}.zip"
spath="~/phd/projfpcc/restserver/logs"
ssh -i ~/fireman.sururu.key ubuntu@${SERVER_IP} "cd ${spath}; zip ${sfname} *.csv"
scp -i ~/fireman.sururu.key ubuntu@${SERVER_IP}:${spath}/${sfname} ~/expresults_fpcc/
scp ~/expresults_fpcc/${sfname} danielfireman@siri.lsd.ufcg.edu.br:/tmp/fireman_${sfname}
