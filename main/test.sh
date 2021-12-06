#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8001 -api=1 -apiPort=9998 &
./server -port=8002 -api=1 -apiPort=9999 &
./server -port=8003 &

sleep 1
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9998/api?key=Tom" &
curl "http://localhost:9998/api?key=Tom" &
curl "http://localhost:9998/api?key=Tom" &

wait