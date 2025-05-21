package api_test

import (
	"fmt"
	"testing"

	"sgridnext.com/server/SgridNodeServer/api"
	"sgridnext.com/src/constant"
)

func TestGetFile(t *testing.T) {
	api.GetFile(api.FileReq{
		ServerId: 4,
		FileName: "SgridTestJavaServer.tar_6758c08789ceb20d5b23eefb5714faddb9bca59318c1c0b28f54a4b40999fc37.gz",
		Type:     constant.FILE_TYPE_PACKAGE,
	})
}

func TestGetConfigList(t *testing.T) {
	fmt.Println("TestGetConfigList INIT")
	api.GetConfigList(api.GetConfigListReq{
		ServerId: 4,
	})
}
