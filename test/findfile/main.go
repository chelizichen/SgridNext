package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)


func FindFile(targetDir string, suffix string) (string, error) {
	var foundFile string
	fmt.Printf("开始在目录下查找文件: %s \n", targetDir)
	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, suffix) {
			foundFile = path
			return filepath.SkipDir
		}
		return nil
	})

	if foundFile == "" {
		return "", fmt.Errorf("未找到%s后缀文件", suffix)
	}
	return foundFile, err
}

func main(){
    file,err := FindFile("/Users/leemulus/Desktop/Boelink/RuoYi-Vue/ruoyi-admin/target",".jar")
    if err != nil{
        fmt.Println("错误:",err.Error())
        return
    }
    fmt.Println("file >> ",file)

}