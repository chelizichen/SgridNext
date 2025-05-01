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
	NodeId     int
}

func (s *ServerInfo) CreateCommand()( *command.Command,error) {
	var cmd *command.Command
	var err error
	if s.ServerType == constant.SERVER_TYPE_BINARY {
		cmd,err= command.CreateBinaryCommand(s.ServerName, s.TargetFile)
	}
	if s.ServerType == constant.SERVER_TYPE_NODE {
		cmd,err = command.CreateNodeCommand(s.ServerName, s.TargetFile)
	}
	if s.ServerType == constant.SERVER_TYPE_JAVA {
		cmd,err = command.CreateJavaJarCommand(s.ServerName, s.TargetFile)
	}
	if err!= nil {
		return nil,err
	}
	cmd.SetNodeId(s.NodeId)
	return cmd,nil
}
