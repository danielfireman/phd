#!/bin/bash
set -x
killall main
for i in `seq 1 $1`
do
	go run main.go > ~/logs/client_$i &
done
