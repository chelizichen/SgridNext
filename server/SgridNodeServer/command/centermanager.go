package command

var CenterManager = NewCenterManager()

type centerManager struct {
	commands map[int]*Command
}

func NewCenterManager() *centerManager {
	return &centerManager{
		commands: make(map[int]*Command),
	}
}
func (cm *centerManager) AddCommand(nodeId int, cmd *Command) {
	cm.commands[nodeId] = cmd
}
func (cm *centerManager) RemoveCommand(nodeId int) {
	delete(cm.commands, nodeId)
}
func (cm *centerManager) GetCommand(nodeId int) *Command {
	return cm.commands[nodeId]
}

