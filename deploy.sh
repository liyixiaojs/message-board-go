#!/bin/sh
# rm -rf ./main;
# GOARCH=amd64 GOOS=linux go build -o main *.go;

echo 'Step 1'
CP="ls ~/golang;"
ssh root@182.92.79.110 "${CP}"

echo 'Done'