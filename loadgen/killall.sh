#!/bin/bash
source configrc

echo "Time to kill bastard loadgen processes running on ${CLIENTS[@]}"
for CLIENT in ${CLIENTS[@]};
do
	ssh -i ~/fireman.sururu.key ubuntu@${CLIENT} 'killall loadgen bash > /dev/null 2>/dev/null'
	echo "Loadgen process at ${CLIENT} ... killed."
done
echo "All killed. Please give me a harder task next time."
