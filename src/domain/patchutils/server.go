package patchutils

import (
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/command"
)

type ServerInfo struct {
	ServerType int
	ServerName string
	BindPort   int
	BindHost   string
	TargetFile string
}

func (s *ServerInfo) CreateCommand()( *command.Command,error) {
	if s.ServerType == constant.SERVER_TYPE_BINARY {
		return command.CreateBinaryCommand(s.ServerName, s.TargetFile)
	}
	if s.ServerType == constant.SERVER_TYPE_NODE {
		return command.CreateNodeCommand(s.ServerName, s.TargetFile)
	}
	if s.ServerType == constant.SERVER_TYPE_JAVA {
		return command.CreateJavaJarCommand(s.ServerName, s.TargetFile)
	}
	return nil,nil
}
