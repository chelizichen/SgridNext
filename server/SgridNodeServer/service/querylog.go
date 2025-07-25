package service

import (
	"sgridnext.com/src/constant"
)

// tail -50000 ${logFile}|tail -${len} | iconv -c -f UTF-8 -t UTF-8|sed 's/[\cA-\cZ]//g'
// tail -50000 ${logFile} | grep -a ${keyword}|tail -${len} | iconv -c -f UTF-8 -t UTF-8|sed 's/[\cA-\cZ]//g'
func QueryLog(logFile string, logType int32, keyword string, len int32) ([]string, error) {
	return constant.QueryLog(logFile, logType, keyword, len)
}

// head -500000 /Users/leemulus/Desktop/临时/waterfull.log |head -a '东兴金蟾'|tail -100 | iconv -c -f UTF-8 -t UTF-8|sed 's/[\cA-\cZ]//g'
// func test() {
// 	cwd, _ := os.Getwd()
// 	fp := filepath.Join(cwd, "waterfull.log")
// 	logRsp, err := searchLog(fp, HEAD, "2025-01-15 15:12:53", 100)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println("logRsp >> ", logRsp)
// }
