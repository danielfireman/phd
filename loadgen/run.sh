#!/bin/bash
set -x
killall main
for i in `seq 1 $1`
do
	sleep 30
	go run main.go > logs/client_${i}_1 &
	go run main.go > logs/client_${i}_2
	sleep 30
done
