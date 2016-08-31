#!/bin/bash
./killall.sh

source configrc

echo "Syncing clients: ${CLIENTS[@]}"
cd ~/phd/loadgen && go build main.go && echo "Loadgen binary compiled"
for CLIENT in ${CLIENTS[@]};
do
	ssh -i ~/fireman.sururu.key ubuntu@${CLIENT} 'cd ~/phd; git pull' &&  echo "Git repository at ${CLIENT} ... synced."
	scp -i ~/fireman.sururu.key main ubuntu@${CLIENT}:~/phd/loadgen/loadgen  &&  echo "Loadgen binary at ${CLIENT} ... synced."
done
echo "All synced. Please give me a harder task next time."
