# rm -rf ./main;
# GOARCH=amd64 GOOS=linux go build -o main *.go;
scp ./main root@182.92.79.110:~/golang/;
scp ./conf.yml root@182.92.79.110:~/golang/;
CP="cd ~/golang;";
ssh root@182.92.79.110 "${CP}";