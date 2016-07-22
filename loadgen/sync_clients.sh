#!/bin/bash
./killall.sh

source configrc

echo "Syncing clients: ${CLIENTS[@]}"
for CLIENT in ${CLIENTS[@]};
do
	ssh -i ~/fireman.sururu.key ubuntu@${CLIENT} 'cd ~/phd; git pull' &&  echo "Git repository at ${CLIENT} ... synced."
	scp -i ~/fireman.sururu.key ~/loadgen ubuntu@${CLIENT}:~/phd/loadgen  &&  echo "Loadgen binary at ${CLIENT} ... synced."
done
echo "All synced. Please give me a harder task next time."
