#!/bin/bash
set -x
rm logs/*
for i in `seq 1 $1`
do
	export CPU_LOG_FILE=cpu_$i; mvn jooby:run
	killall java
	sleep 10s
done
