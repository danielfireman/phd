#!/bin/bash
set -x

declare -a CLIENTS=('localhost' 'localhost')
NUM_ROUNDS=30
SERVER=http://localhost:8080

for CLIENT in ${CLIENTS[@]};
do
	ssh ${CLIENT} 'cd phd/loadgen;./run ${NUM_ROUNDS} ${SERVER}' &
done
