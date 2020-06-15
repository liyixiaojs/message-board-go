#!/bin/sh
# rm -rf ./main;
# GOARCH=amd64 GOOS=linux go build -o main *.go;
ssh-keyscan -H 182.92.79.110 >> ~/.ssh/known_hosts
echo 'Step 1'
scp -o StrictHostKeyChecking=no ./main root@182.92.79.110:~/golang/
echo 'Step 2'
scp -o StrictHostKeyChecking=no ./conf.yaml root@182.92.79.110:~/golang/

echo 'Step 3'
CP="cd ~/golang;"
ssh root@182.92.79.110 "${CP}"

echo 'Done'