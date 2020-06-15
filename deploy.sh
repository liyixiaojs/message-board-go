#!/bin/sh
# rm -rf ./main;
# GOARCH=amd64 GOOS=linux go build -o main *.go;
echo 'Step 1'
scp ./main root@182.92.79.110:~/golang/;
echo 'Step 2'
scp ./conf.yaml root@182.92.79.110:~/golang/;

echo 'Step 3'
CP="cd ~/golang;";
ssh root@182.92.79.110 "${CP}";

echo 'Done'