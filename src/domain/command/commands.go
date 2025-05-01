package command

func CreateNodeCommand(serverName string, targetFile string) (*Command ,error){
	cmd := NewServerCommand(serverName)
	err := cmd.SetCommand("node", targetFile)
	return cmd,err
}

func CreateBinaryCommand(serverName string, targetFile string) (*Command ,error){
	cmd := NewServerCommand(serverName)
	err := cmd.SetCommand(targetFile)
	return cmd,err
}

func CreateJavaJarCommand(serverName string, targetFile string) (*Command ,error){
	cmd := NewServerCommand(serverName)
	err := cmd.SetCommand("java", "-jar", targetFile)
	return cmd,err
}
