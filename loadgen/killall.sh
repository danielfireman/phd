#!/bin/bash
declare -a CLIENTS=('10.4.5.132' '10.4.5.130' '10.4.5.133' '10.4.5.134')

echo "Time to kill bastard loadgen processes running on ${CLIENTS[@]}"
for CLIENT in ${CLIENTS[@]};
do
	ssh -i ~/fireman.sururu.key ubuntu@${CLIENT} 'killall loadgen > /dev/null 2>/dev/null'
	echo "Loadgen process at ${CLIENT} ... killed."
done
echo "All killed. Please give me a harder task next time."
