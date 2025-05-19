#!/bin/bash

# if permission denied
# run script with ` chmod +x build.sh `
readonly ServerName="app"

rm -r dist
cd web 
npm run build
cd ..
cp -r web/dist dist

# compile
# GOOS=linux GOARCH=amd64
rm -r $ServerName
GOOS=linux GOARCH=amd64 go build -o $ServerName

