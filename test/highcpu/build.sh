#!/bin/bash

# if permission denied
# run script with ` chmod +x build.sh `
readonly ServerName="TestHighCpuServer"


# compile
# GOOS=linux GOARCH=amd64
rm -r $ServerName
GOOS=linux GOARCH=amd64 go build -o $ServerName

# build
tar -cvf $ServerName.tar.gz  $ServerName

rm -r $ServerName