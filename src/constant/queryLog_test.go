package constant

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestQueryLog(t *testing.T) {
	// 创建临时测试文件
	testContent := `line1: error occurred
line2: warning message
line3: info message
line4: error and warning
line5: another error
line6: final warning`

	tmpfile, err := ioutil.TempFile("", "test_log_*.txt")
	if err != nil {
		t.Fatalf("无法创建临时文件: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(testContent)); err != nil {
		t.Fatalf("无法写入临时文件: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("无法关闭临时文件: %v", err)
	}

	// 测试用例
	tests := []struct {
		name      string
		logType   int32
		keyword   string
		lineCount int32
		expected  int // 期望的结果行数
	}{
		{"无关键词", TAIL, "", 6, 6},
		{"单个关键词", TAIL, "error", 6, 3},
		{"多个关键词 AND", TAIL, "error+warning", 6, 1},
		{"多个关键词带空格", TAIL, "error + warning", 6, 1},
		{"特殊字符", TAIL, "error+.*", 6, 0}, // 应该转义.*
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := QueryLog(tmpfile.Name(), tc.logType, tc.keyword, tc.lineCount)
			if err != nil {
				t.Fatalf("QueryLog失败: %v", err)
			}

			fmt.Printf("测试用例 '%s' 结果: %v\n", tc.name, result)
			if len(result) != tc.expected {
				t.Errorf("期望%d行结果，但得到%d行\n结果: %v",
					tc.expected, len(result), strings.Join(result, "\n"))
			}
		})
	}
}
