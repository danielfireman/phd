#!/bin/bash
set -x

declare -a CLIENTS=('10.4.5.132' '10.4.5.130' '10.4.5.133' '10.4.5.134')
NUM_ROUNDS=30
SERVER=http://10.4.2.103:8080

for CLIENT in ${CLIENTS[@]};
do
	ssh ${CLIENT} 'cd phd/loadgen;./run ${NUM_ROUNDS} ${SERVER}' &
done
