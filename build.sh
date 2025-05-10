#!/bin/bash

# if permission denied
# run script with ` chmod +x build.sh `
readonly ServerName="app"

rm -r dist
mkdir -p dist

# compile
# GOOS=linux GOARCH=amd64
rm -r $ServerName
GOOS=linux GOARCH=amd64 go build -o $ServerName

# build
# tar -cvf $ServerName.tar.gz  $ServerName

# mv $ServerName.tar.gz  dist
# mv $ServerName dist

#  GOOS=linux GOARCH=amd64 go build -o GoTradeBackServer
