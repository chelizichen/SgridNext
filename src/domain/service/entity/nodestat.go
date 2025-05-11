package entity

type NodeStatFactory struct {
	Base *NodeStat
}

func NewNodeStatFactory(Base *NodeStat) *NodeStatFactory {
	return &NodeStatFactory{
		Base: Base,
	}
}

func (f *NodeStatFactory) Assign(NewInstance *NodeStat) *NodeStat {
	if f.Base.Id != 0 {
		NewInstance.Id = f.Base.Id
	}
	if f.Base.NodeId != 0 {
		NewInstance.NodeId = f.Base.NodeId
	}
	if f.Base.ServerNodeId != 0 {
		NewInstance.ServerNodeId = f.Base.ServerNodeId
	}
	if f.Base.ServerName != "" {
		NewInstance.ServerName = f.Base.ServerName
	}
	if f.Base.ServerId != 0 {
		NewInstance.ServerId = f.Base.ServerId
	}
	if f.Base.TYPE != 0 {
		NewInstance.TYPE = f.Base.TYPE
	}
	if f.Base.Content != "" {
		NewInstance.Content = f.Base.Content
	}
	return NewInstance
}
