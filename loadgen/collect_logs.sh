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
	port=$(getport ${CLIENT})
	fname="c${NUM_CORES}_${EXECUTION_ID}_${ACTIVE_CLIENT}.zip"
	ssh -i ~/fireman.sururu.key ${SSH_ADDR} -p ${port}  "cd ~/phd/loadgen/logs; zip ${fname} client*"
	scp -P ${port} -i ~/fireman.sururu.key ${SSH_ADDR}:~/phd/loadgen/logs/${fname} /tmp/explogs_${fname}
	if [ $ACTIVE_CLIENT == $NUM_CLIENTS ]; then
		break
	fi
	ACTIVE_CLIENT=`expr $ACTIVE_CLIENT + 1`
done

echo "Packing server logs"
sfname="s${NUM_CORES}_${EXECUTION_ID}.zip"
spath="~/phd/projfpcc/restserver/logs"
sport=$(getport ${SERVER_IP})
ssh -i ~/fireman.sururu.key ${SSH_ADDR} -p ${sport} "cd ${spath}; zip ${sfname} *.csv"
scp -P ${sport} -i ~/fireman.sururu.key ${SSH_ADDR}:${spath}/${sfname} /tmp/explogs_${sfname}
