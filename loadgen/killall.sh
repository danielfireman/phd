#!/bin/bash
source configrc

echo "Time to kill bastard loadgen processes running on ${CLIENTS[@]}"
for CLIENT in ${CLIENTS[@]};
do
	port=$(getport ${CLIENT})
	echo "Server ${SSH_ADDR}:${port}"
	ssh -i ~/fireman.sururu.key ${SSH_ADDR} -p ${port} 'killall loadgen > /dev/null 2>/dev/null'
	echo "Loadgen process at ${SSH_ADDR}:${port} ... killed."
done
echo "All killed. Please give me a harder task next time."
