readonly ServerName="sgridnode"


# compile
# GOOS=linux GOARCH=amd64
rm -r $ServerName
GOOS=linux GOARCH=amd64 go build -o $ServerName
