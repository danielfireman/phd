#!/bin/bash
set -x
killall main
go run main.go > ~/logs/client_$1 &
