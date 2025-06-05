#!/bin/bash

# if permission denied
# run script with ` chmod +x build.sh `
readonly ServerName="TestBashGoServer"


# compile
# GOOS=linux GOARCH=amd64
rm -r $ServerName
GOOS=linux GOARCH=amd64 go build -o $ServerName

# build
tar -cvf $ServerName.tar.gz  $ServerName ./bashgo.sh

rm -r $ServerName