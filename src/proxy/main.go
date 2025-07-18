package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"sgridnext.com/server/SgridNodeServer/command"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/domain/service/mapper"
	"sgridnext.com/src/logger"
	"sgridnext.com/src/patchutils"
)

type T_Proxy struct {
	NodeId int
	Proxy  *protocol.NodeServantClient
}

type T_ProxyMap struct {
	sync.RWMutex
	items map[int]*T_Proxy
	hostMap map[string]int
}

func (p *T_ProxyMap) AddProxy(nodeId int,host string, proxy *protocol.NodeServantClient) {
	p.Lock()
	defer p.Unlock()
	p.items[nodeId] = &T_Proxy{
		NodeId: nodeId,
		Proxy:  proxy,
	}
	p.hostMap[host] = nodeId
}

// FullDispatch 全量节点调用
func (p *T_ProxyMap) FullDispatch(callback func(*protocol.NodeServantClient) error) ([]int, []int) {
	p.RLock()
	defer p.RUnlock()

	var successIDs, failIDs []int
	for id, proxy := range p.items {
		err := callback(proxy.Proxy)
		if err != nil {
			failIDs = append(failIDs, id)
		} else {
			successIDs = append(successIDs, id)
		}
	}
	return successIDs, failIDs
}

func (p *T_ProxyMap) GetProxyByHost(host string) *protocol.NodeServantClient {
	p.RLock()
	defer p.RUnlock()
	return p.items[p.hostMap[host]].Proxy
}

func (p *T_ProxyMap) DispatchByHost(host string, callback func(*protocol.NodeServantClient) error) error {
	p.RLock()
	defer p.RUnlock()
	proxy := p.items[p.hostMap[host]].Proxy
	if proxy == nil {
		return fmt.Errorf("no available nodes")
	}
	return callback(proxy)
}

// RandomDispatch 随机节点调用
func (p *T_ProxyMap) RandomDispatch(callback func(*protocol.NodeServantClient) error) (int, error) {
	p.RLock()
	defer p.RUnlock()

	if len(p.items) == 0 {
		return -1, fmt.Errorf("no available nodes")
	}
	rand.Seed(time.Now().UnixNano())
	nodes := make([]int, 0, len(p.items))
	for id := range p.items {
		nodes = append(nodes, id)
	}

	selected := nodes[rand.Intn(len(nodes))]
	err := callback(p.items[selected].Proxy)
	return selected, err
}

func (p *T_ProxyMap) RemoveProxy(nodeId int) {
	delete(p.items, nodeId)
}

func (p *T_ProxyMap) GetNodes() []int {
	nodes := make([]int, 0, len(p.items))
	for id := range p.items {
		nodes = append(nodes, id)
	}
	return nodes
}

var ProxyMap *T_ProxyMap

func init() {
	ProxyMap = &T_ProxyMap{
		items: make(map[int]*T_Proxy),
		hostMap: make(map[string]int),
		RWMutex: sync.RWMutex{},
	}
}

func LoadProxy() {
	// 定时同步节点信息
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		NodesStatMap := &command.SvrNodeStatMap{
			UpdateTime: constant.GetCurrentTime(),
			StatList: make([]*command.SvrNodeStat, 0),
		}
		for range ticker.C {
			NodesStatMap.UpdateTime = constant.GetCurrentTime()
			NodesStatMap.StatList = nil
			NodesStatMap.StatList = make([]*command.SvrNodeStat, 0)
			// 拉取所有 节点
			nodes, err := mapper.T_Mapper.GetNodeList()
			if err != nil {
				logger.App.Errorf("同步节点失败 | %s", err.Error())
				continue
			}
			for _, node := range nodes {
				if node.NodeStatus == constant.COMM_STATUS_ONLINE {
					if patchutils.T_PatchUtils.Contains(ProxyMap.GetNodes(), node.ID) {
						// 在线,并且也在节点中
						_, err := (*ProxyMap.items[node.ID].Proxy).KeepAlive(context.Background(), &emptypb.Empty{})
						if err != nil {
							mapper.T_Mapper.UpdateMachineNodeStatus(node.ID, constant.COMM_STATUS_OFFLINE)
							logger.Alive.Errorf("节点 | %s | 挂了 | %s", node.ID, err.Error())
							continue
						}
						nodeStatData,err  := (*ProxyMap.items[node.ID].Proxy).GetNodeStat(context.Background(), &emptypb.Empty{})
						if err != nil{
							logger.Alive.Errorf("节点 | %s | 获取状态异常 | %s", node.ID, err.Error())
						}
						var svrNodeMap *command.SvrNodeStatMap
						json.Unmarshal([]byte(nodeStatData.Data),&svrNodeMap)
						NodesStatMap.StatList = append(NodesStatMap.StatList, svrNodeMap.StatList...)
						mapper.T_Mapper.UpdateMachineNodeStatus(node.ID, constant.COMM_STATUS_ONLINE)
					} else {
						// 在线，不在节点中，可能为新添加的
						addr := fmt.Sprintf("%s:%s", node.Host, constant.NODE_PORT)
						conn, err := grpc.NewClient(addr,
							grpc.WithTransportCredentials(
								insecure.NewCredentials(),
							),
							grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(constant.MSG_RECV_SIZE_MAX)),
							grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(constant.MSG_CALL_SIZE_MAX)),
						)
						if err != nil {
							mapper.T_Mapper.UpdateMachineNodeStatus(node.ID, constant.COMM_STATUS_OFFLINE)
							logger.App.Errorf("创建节点连接失败 | ID:%d | %s", node.ID, err.Error())
							continue
						}
						client := protocol.NewNodeServantClient(conn)
						ProxyMap.AddProxy(node.ID, node.Host, &client)
						logger.Alive.Infof("添加节点成功 | %s ", node.ID)
						mapper.T_Mapper.UpdateMachineNodeStatus(node.ID, constant.COMM_STATUS_ONLINE)
					}
				}
				if node.NodeStatus == constant.COMM_STATUS_OFFLINE {
					if patchutils.T_PatchUtils.Contains(ProxyMap.GetNodes(), node.ID) {
						// 下线了，仍然在 节点中
						ProxyMap.RemoveProxy(node.ID)
					} else {
						// pass
					}
				}
			}
			jsonStr,err  := json.Marshal(NodesStatMap)
			for _, node := range nodes{
				if node.NodeStatus != constant.COMM_STATUS_ONLINE {
					continue
				}
				if !patchutils.T_PatchUtils.Contains(ProxyMap.GetNodes(), node.ID) {
					continue
				}
				if err != nil {
					logger.Alive.Errorf("全量节点同步失败 | 序列化 ｜ %s",err.Error())
					continue
				}
				syncRsp,err := (*ProxyMap.items[node.ID].Proxy).SyncAllNodeStat(context.Background(),&protocol.SyncStatReq{
					Data: string(jsonStr),
				})
				if err != nil{
					logger.Alive.Errorf("全量节点同步失败 | 发送失败 ｜ %s",err.Error())
					continue
				}
				logger.Alive.Infof("全量同步成功 | %v",syncRsp)
			}
			logger.Alive.Info("全量同步完成,开始同步到本地节点")
			// 本地备份节点状态
			cwd, _ := os.Getwd()
			stat_remote_path := filepath.Join(cwd, "stat-remote.json")
			logger.Alive.Infof("开始同步到本地节点路径 ｜ %s",stat_remote_path)
			outFile, err := os.OpenFile(stat_remote_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				logger.App.Errorf("创建文件失败: SyncStat |%v", err)
			}
			defer outFile.Close()
			if _, err := outFile.Write([]byte(jsonStr)); err != nil {
				logger.App.Errorf("文件写入失败: SyncStat | %v", err)
			}
		}
	}()
	// ticker.Stop()
}
