#!/bin/bash
trap "rm server;kill 0" EXIT

rm -rf ./logs/n1/*

go build -o server
./server -c ./config/n1_master.yaml &
./server -c ./config/n1_slaver_1.yaml &
./server -c ./config/n1_slaver_2.yaml
#./server -c ./config/n2_master.yaml &
#./server -c ./config/n2_slaver_1.yaml &
#./server -c ./config/n2_slaver_2.yaml &
#./server -c ./config/n3_master.yaml &
#./server -c ./config/n3_slaver_1.yaml &
#./server -c ./config/n3_slaver_2.yaml
wait