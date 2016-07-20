#!/bin/bash
declare -a CLIENTS=('10.4.5.132' '10.4.5.130' '10.4.5.133' '10.4.5.134')

echo "Syncing clients: ${CLIENTS[@]}"
for CLIENT in ${CLIENTS[@]};
do
	ssh -i ~/fireman.sururu.key ubuntu@${CLIENT} 'cd ~/phd; git pull' &&  echo "Git repository at at ${CLIENT} ... synced."
done
echo "All synced. Please give me a harder task next time."
