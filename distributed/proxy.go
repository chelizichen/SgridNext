package distributed

import (
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sgridnext.com/server/SgridNodeServer/command"
)

type PrxManage[T any] struct {
	svr       string
	addrs     []*command.BaseSvrNodeStat
	clients   []*T
	newClient func(conn grpc.ClientConnInterface) T
	getAddrs  func() []*command.BaseSvrNodeStat
}

type Proxy[T any] interface {
	GetAddrs() []*command.BaseSvrNodeStat
	NewClient(conn grpc.ClientConnInterface) T
}

func LoadStringToProxy[T any](prx Proxy[T]) (*PrxManage[T], error) {
	addrs := prx.GetAddrs()
	if len(addrs) == 0 {
		panic("no registry found")
	}
	first := addrs[0]
	pm := &PrxManage[T]{
		addrs:     addrs,
		svr:       first.ServerName,
		getAddrs:  prx.GetAddrs,
		newClient: prx.NewClient,
		clients:   make([]*T, 0),
	}
	pm.syncNodes(true)
	go pm.syncNodes(false)
	return pm, nil
}

func (p *PrxManage[T]) GetClient() T {
	if len(p.clients) == 0 {
		var zero T
		return zero
	}
	idx := rand.Intn(len(p.clients))
	return *p.clients[idx]
}

func (p *PrxManage[T]) syncNodes(init bool) {
	if init {
		p.addrs = p.getAddrs()
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
			}
			if len(diffSub) > 0 {
				err := p.delNodes(diffSub)
				if err != nil {
					fmt.Println("同步删除节点失败 | ", err.Error())
				}
			}
		}
	}
}

func (p *PrxManage[T]) addNodes(addrs []*command.BaseSvrNodeStat) error {
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
		p.clients = append(p.clients, &cn)
	}
	return nil
}

func (p *PrxManage[T]) delNode(addr *command.BaseSvrNodeStat) {
	for i, client := range p.clients {
		if client == nil {
			continue
		}
		p.clients[i] = nil
	}
}

func (p *PrxManage[T]) delNodes(addrs []*command.BaseSvrNodeStat) error {
	if len(addrs) == 0 {
		return nil
	}
	for _, addr := range addrs {
		p.delNode(addr)
	}
	return nil
}

func getAddrAdded(oldAddrs, newAddrs []*command.BaseSvrNodeStat) []*command.BaseSvrNodeStat {
	added := []*command.BaseSvrNodeStat{}
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

func getAddrSubed(oldAddrs, newAddrs []*command.BaseSvrNodeStat) []*command.BaseSvrNodeStat {
	subed := []*command.BaseSvrNodeStat{}
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
