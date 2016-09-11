#!/bin/bash

./killall.sh

source configrc

echo "Syncing clients: ${CLIENTS[@]}"
cd ~/phd/loadgen && go build main.go && echo "Loadgen binary compiled"
for CLIENT in ${CLIENTS[@]};
do
	port=$(getport ${CLIENT})
	echo "Server ${SSH_ADDR}:${port}"
	ssh -i ~/fireman.sururu.key ${SSH_ADDR} -p ${port} 'cd ~/phd; git pull' &&  echo "Git repository at port ${SSH_ADDR}:${port} ... synced."
	scp -P ${port} -i ~/fireman.sururu.key main ${SSH_ADDR}:~/phd/loadgen/loadgen &&  echo "Loadgen binary at ${SSH_ADDR}:${port} ... synced."
done
echo "All synced. Please give me a harder task next time."
