package command

import "sgridnext.com/src/constant"

type ServerInfo struct {
	ServerType int
	ServerName string
	BindPort   int
	BindHost   string
	TargetFile string
	NodeId     int
}

func (s *ServerInfo) CreateCommand() (*Command, error) {
	var cmd *Command
	var err error
	if s.ServerType == constant.SERVER_TYPE_BINARY {
		cmd, err = CreateBinaryCommand(s.ServerName, s.TargetFile)
	}
	if s.ServerType == constant.SERVER_TYPE_NODE {
		cmd, err = CreateNodeCommand(s.ServerName, s.TargetFile)
	}
	if s.ServerType == constant.SERVER_TYPE_JAVA {
		cmd, err = CreateJavaJarCommand(s.ServerName, s.TargetFile)
	}
	if err != nil {
		return nil, err
	}
	cmd.SetNodeId(s.NodeId)
	return cmd, nil
}
