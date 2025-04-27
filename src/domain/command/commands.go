package command

func CreateNodeCommand(serverName string, targetFile string) *Command {
	cmd := NewServerCommand(serverName)
	cmd.SetCommand("node", targetFile)
	return cmd
}

func CreateBinaryCommand(serverName string, targetFile string) *Command {
	cmd := NewServerCommand(serverName)
	cmd.SetCommand(targetFile)
	return cmd
}

func CreateJavaJarCommand(serverName string, targetFile string) *Command {
	cmd := NewServerCommand(serverName)
	cmd.SetCommand("java", "-jar", targetFile)
	return cmd
}
