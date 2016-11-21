#!/bin/bash
set -x
rm logs/*
mvn clean package -DskipTests
for i in `seq 1 $1`
do
	VM_FLAGS=""
	export ROUND=$i; export CPU_LOG_FILE=cpu_$i; 
	if [ -n "${CONTROL_GC}" ]; then
		VM_FLAGS="-Xmx1024m -Xms1024m -XX:NewRatio=1 -XX:InitialCodeCacheSize=32m -XX:MetaspaceSize=32m"
	fi
	java ${VM_FLAGS} -jar target/restserver-1.0.jar
	killall java
	sleep 10s
done
