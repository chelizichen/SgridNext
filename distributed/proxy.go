package distributed

import (
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
type BaseSvrNodeStat struct {
	ServerName string `json:"server_name,omitempty"` // 服务名称
	ServerHost string `json:"host,omitempty"`        // 主机地址
	ServerPort int    `json:"port,omitempty"`        // 主机端口
}

func (p *BaseSvrNodeStat) String() string {
	return fmt.Sprintf("%s:%d", p.ServerHost, p.ServerPort)
}

type PrxManage[T any] struct {
	svr       string
	addrs     []*BaseSvrNodeStat
	clients   map[string]*T
	newClient func(conn grpc.ClientConnInterface) T
	getAddrs  func() []*BaseSvrNodeStat
}

type Proxy[T any] interface {
	GetAddrs() []*BaseSvrNodeStat
	NewClient(conn grpc.ClientConnInterface) T
	GetServerName() string
}

func LoadStringToProxy[T any](prx Proxy[T]) (*PrxManage[T], error) {
	pm := &PrxManage[T]{
		addrs:     make([]*BaseSvrNodeStat, 0),
		svr:       prx.GetServerName(),
		getAddrs:  prx.GetAddrs,
		newClient: prx.NewClient,
		clients:   make(map[string]*T, 0),
	}
	pm.syncNodes(true)
	go pm.syncNodes(false)
	return pm, nil
}

func (p *PrxManage[T]) GetClient() (client T,ok bool){
	if len(p.clients) == 0 {
		return client, false
	}
	rand.Seed(time.Now().UnixNano())
	keys := make([]string, 0, len(p.clients))
	for k := range p.clients {
		keys = append(keys, k)
	}
	randomKey := keys[rand.Intn(len(keys))]
	client = *p.clients[randomKey]
	return client, true
}

func (p *PrxManage[T]) syncNodes(init bool) {
	if init {
		p.addrs = p.getAddrs()
		if len(p.addrs) == 0 {
			fmt.Println("初始化同步节点失败 | 节点列表为空")
			return
		}
		err := p.addNodes(p.addrs)
		if err != nil {
			fmt.Println("初始化同步节点失败 | ", err.Error())
		}
		return
	}
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			newAddrs := p.getAddrs()
			oldAddrs := p.addrs
			diffAdd := getAddrAdded(oldAddrs, newAddrs)
			diffSub := getAddrSubed(oldAddrs, newAddrs)
			if len(diffAdd) > 0 {
				err := p.addNodes(diffAdd)
				if err != nil {
					fmt.Println("同步新增节点失败 | ", err.Error())
				}
				fmt.Println("同步新增节点成功 | ", diffAdd)
			}
			if len(diffSub) > 0 {
				err := p.delNodes(diffSub)
				if err != nil {
					fmt.Println("同步删除节点失败 | ", err.Error())
				}
				fmt.Println("同步删除节点成功 | ", diffSub)
			}
			p.addrs = newAddrs
			fmt.Println("同步节点成功 | ", p.clients)
		}
	}
}

func (p *PrxManage[T]) addNodes(addrs []*BaseSvrNodeStat) error {
	for _, addr := range addrs {
		t_addr := fmt.Sprintf("%s:%v", addr.ServerHost, addr.ServerPort)
		conn, err := grpc.NewClient(t_addr,
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		)
		if err != nil {
			fmt.Println("创建节点连接失败 | ", t_addr, err.Error())
			return err
		}
		cn := p.newClient(conn)
		p.clients[addr.String()] = &cn
	}
	return nil
}

func (p *PrxManage[T]) delNode(addr *BaseSvrNodeStat) {
	p.clients[addr.String()] = nil
	delete(p.clients, addr.String())
}

func (p *PrxManage[T]) delNodes(addrs []*BaseSvrNodeStat) error {
	if len(addrs) == 0 {
		return nil
	}
	for _, addr := range addrs {
		p.delNode(addr)
	}
	return nil
}

func getAddrAdded(oldAddrs, newAddrs []*BaseSvrNodeStat) []*BaseSvrNodeStat {
	added := []*BaseSvrNodeStat{}
	oldMap := make(map[string]struct{})
	for _, addr := range oldAddrs {
		key := fmt.Sprintf("%s:%d", addr.ServerHost, addr.ServerPort)
		oldMap[key] = struct{}{}
	}
	for _, addr := range newAddrs {
		key := fmt.Sprintf("%s:%d", addr.ServerHost, addr.ServerPort)
		if _, ok := oldMap[key]; !ok {
			added = append(added, addr)
		}
	}
	return added
}

func getAddrSubed(oldAddrs, newAddrs []*BaseSvrNodeStat) []*BaseSvrNodeStat {
	subed := []*BaseSvrNodeStat{}
	newMap := make(map[string]struct{})
	for _, addr := range newAddrs {
		key := fmt.Sprintf("%s:%d", addr.ServerHost, addr.ServerPort)
		newMap[key] = struct{}{}
	}
	for _, addr := range oldAddrs {
		key := fmt.Sprintf("%s:%d", addr.ServerHost, addr.ServerPort)
		if _, ok := newMap[key]; !ok {
			subed = append(subed, addr)
		}
	}
	return subed
}
