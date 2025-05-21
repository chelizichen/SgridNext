package constant

import "time"

func ConvertToInt32Slice(input []int) []int32 {
	// 创建一个新的 int32 切片
	result := make([]int32, len(input))
	// 遍历 input 切片，将每个元素转换为 int32 并添加到 result 切片中
	for i, v := range input {
		result[i] = int32(v)
	}
	return result
}

func ConvertToIntSlice(input []int32) []int {
	// 创建一个新的 int32 切片
	result := make([]int, len(input))
	// 遍历 input 切片，将每个元素转换为 int32 并添加到 result 切片中
	for i, v := range input {
		result[i] = int(v)
	}
	return result
}

// 拿到当前时间 YYYY-MM-DD HH:MM:SS
func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
