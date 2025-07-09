package command

import (
	"sgridnext.com/src/config"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/logger"
)

type ServerInfo struct {
	ServerType int
	ServerName string
	BindPort   int
	BindHost   string
	TargetFile string
	NodeId     int
	AdditionalArgs string
	ServerRunType int
	ServerId int
	DockerName string
}

func (s *ServerInfo) CreateCommand() (*Command, error) {
	var cmd *Command
	var err error
	if s.ServerType == constant.SERVER_TYPE_BINARY {
		cmd, err = CreateBinaryCommand(s)
	}
	if s.ServerType == constant.SERVER_TYPE_NODE {
		cmd, err = CreateNodeCommand(s)
	}
	if s.ServerType == constant.SERVER_TYPE_JAVA {
		cmd, err = CreateJavaJarCommand(s)
	}
	if err != nil {
		return nil, err
	}
	localNodeId := config.Conf.GetLocalNodeId()
	cmd.SetNodeId(s.NodeId)
	cmd.SetHost(s.BindHost)
	cmd.SetPort(s.BindPort)
	logger.CMD.Infof("debug >>> CreateCommand | SetHost %s | SetPort %d", s.BindHost, s.BindPort)
	cmd.SetLocalMachineId(localNodeId)
	cmd.SetServerId(s.ServerId)
	cmd.SetAdditionalArgs(s.AdditionalArgs)
	cmd.SetRunServerType(s.ServerRunType)
	cmd.SetDockerName(s.DockerName)
	return cmd, nil
}
