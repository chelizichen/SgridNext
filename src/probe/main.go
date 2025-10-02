package probe

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	protocol "sgridnext.com/server/SgridNodeServer/proto"
	"sgridnext.com/src/constant"
	"sgridnext.com/src/logger"
)

// 进行 SgridServer 的 服务探测
// 端口是 25528 ，通过输入 10.122.**.*
// 将 10.122 前缀IP的 所有25528端口进行探测，调用 SgridServer 的 Probe 接口
// 如果能够连通，则返回 IP地址，最终返回成一个列表

var (
	timeout   = 5 * time.Second
	maxIPs    = 0
	batchSize = 500
)

type ProbeResult struct {
	IP     string
	Status string
	Error  string
}

func Probe(networkPrefixs []string) []ProbeResult {
	logger.Probe.Infof("开始探测网段: %v\n", networkPrefixs)
	logger.Probe.Infof("探测端口: %s\n", constant.NODE_PORT)
	logger.Probe.Infof("超时时间: %v\n", timeout)
	var results []ProbeResult = make([]ProbeResult, 0)
	for _, networkPrefix := range networkPrefixs {
		// 验证网段前缀格式
		if _, err := parseNetworkPrefix(networkPrefix); err != nil {
			logger.Probe.Infof("错误: %v\n", err)
		}

		startTime := time.Now()
		results = append(results, probeNetwork(networkPrefix)...)
		duration := time.Since(startTime)

		successCount := 0
		portOpenCount := 0

		for _, result := range results {
			if result.Status == "成功" {
				successCount++
			} else if result.Status == "端口不可达" {
				// 不显示端口不可达的详细信息，减少输出噪音
				portOpenCount++
			}
		}

		logger.Probe.Infof("\n统计信息:\n")
		logger.Probe.Infof("- 总扫描IP数: %d\n", len(results))
		logger.Probe.Infof("- 端口开放数: %d\n", portOpenCount)
		logger.Probe.Infof("- 服务可用数: %d\n", successCount)
		logger.Probe.Infof("- 扫描耗时: %v\n", duration)

		if successCount > 0 {
			logger.Probe.Infof("\n🎉 发现 %d 个可用的 SgridServer 节点\n", successCount)
		} else {
			logger.Probe.Infof("\n⚠️  未发现可用的 SgridServer 节点")
		}
	}
	return results
}

// 探测指定网段的所有IP
func probeNetwork(networkPrefix string) []ProbeResult {
	var allResults []ProbeResult

	// 解析网段前缀，确定扫描范围
	ipList := generateIPList(networkPrefix)

	// 如果设置了最大IP数量限制
	if maxIPs > 0 && len(ipList) > maxIPs {
		logger.Probe.Infof("警告: 网段包含 %d 个IP，超过限制 %d，将只扫描前 %d 个IP\n",
			len(ipList), maxIPs, maxIPs)
		ipList = ipList[:maxIPs]
	}

	logger.Probe.Infof("将扫描 %d 个IP地址，批量大小: %d\n", len(ipList), batchSize)

	// 分批扫描
	for i := 0; i < len(ipList); i += batchSize {
		end := i + batchSize
		if end > len(ipList) {
			end = len(ipList)
		}

		batch := ipList[i:end]
		logger.Probe.Infof("扫描批次 %d-%d: %v\n", i+1, end, batch)

		batchResults := probeBatch(batch)
		allResults = append(allResults, batchResults...)

		// 如果不是最后一批，稍作停顿
		if end < len(ipList) {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return allResults
}

// 批量扫描IP
func probeBatch(ipList []string) []ProbeResult {
	var results []ProbeResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, ip := range ipList {
		wg.Add(1)

		go func(targetIP string) {
			defer wg.Done()
			result := probeSingleHost(targetIP)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(ip)
	}

	wg.Wait()
	return results
}

// 生成IP地址列表
func generateIPList(networkPrefix string) []string {
	var ipList []string
	parts := strings.Split(networkPrefix, ".")

	if len(parts) == 2 {
		// /16 网段: 10.169 -> 10.169.0.1 到 10.169.255.254
		logger.Probe.Infof("扫描 /16 网段: %s.*.*\n", networkPrefix)
		for i := 0; i <= 255; i++ {
			for j := 1; j <= 254; j++ {
				ip := fmt.Sprintf("%s.%d.%d", networkPrefix, i, j)
				ipList = append(ipList, ip)
			}
		}
	} else if len(parts) == 3 {
		// /24 网段: 10.169.114 -> 10.169.114.1 到 10.169.114.254
		logger.Probe.Infof("扫描 /24 网段: %s.*\n", networkPrefix)
		for i := 1; i <= 254; i++ {
			ip := fmt.Sprintf("%s.%d", networkPrefix, i)
			ipList = append(ipList, ip)
		}
	} else {
		// 其他情况，默认按 /24 处理
		logger.Probe.Infof("扫描 /24 网段: %s.*\n", networkPrefix)
		for i := 1; i <= 254; i++ {
			ip := fmt.Sprintf("%s.%d", networkPrefix, i)
			ipList = append(ipList, ip)
		}
	}

	return ipList
}

// 探测单个主机
func probeSingleHost(ip string) ProbeResult {
	// 首先检查端口是否开放
	if !isPortOpen(ip, constant.NODE_PORT) {
		return ProbeResult{
			IP:     ip,
			Status: "端口不可达",
			Error:  "连接超时",
		}
	}

	// 尝试建立gRPC连接并调用Probe接口
	address := fmt.Sprintf("%s:%s", ip, constant.NODE_PORT)

	// 创建gRPC连接
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return ProbeResult{
			IP:     ip,
			Status: "gRPC连接失败",
			Error:  err.Error(),
		}
	}
	defer conn.Close()

	// 创建客户端
	client := protocol.NewNodeServantClient(conn)

	// 调用Probe接口
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = client.Probe(ctx, &emptypb.Empty{})
	if err != nil {
		return ProbeResult{
			IP:     ip,
			Status: "Probe调用失败",
			Error:  err.Error(),
		}
	}

	return ProbeResult{
		IP:     ip,
		Status: "成功",
		Error:  "",
	}
}

// 检查端口是否开放
func isPortOpen(host, port string) bool {
	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// 解析网段前缀，支持更灵活的输入
func parseNetworkPrefix(prefix string) (string, error) {
	parts := strings.Split(prefix, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("无效的网段前缀: %s", prefix)
	}

	// 验证每个部分都是数字
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return "", fmt.Errorf("无效的网段前缀: %s", prefix)
		}
	}

	return prefix, nil
}
