package resource

type NodeResource struct {
	SystemInfo  SystemInfo    `json:"systemInfo"`
	ProcessInfo []ProcessInfo `json:"processInfo"`
}

func GetNodeResource(workDir string) NodeResource {
	// 获取系统状态
	systemStatus := GetSystemStats()
	// 获取进程状态
	processStatus := GetProcessInfo(workDir)
	return NodeResource{
		SystemInfo:  systemStatus,
		ProcessInfo: processStatus,
	}
}
