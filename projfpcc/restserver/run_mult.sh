#!/bin/bash
set -x
for i in `seq 1 $1`
do
	export CPU_LOG_FILE=cpu_$i; mvn jooby:run
done