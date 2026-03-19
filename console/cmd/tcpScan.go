package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "net:scan",
		Short: "TCP 端口扫描工具",
		Run:   runTcpScan,
	}
	cmd.Flags().String("ip", "127.0.0.1", "IP address")
	cmd.Flags().Int("port1", 1, "start port")
	cmd.Flags().Int("port2", 65535, "end port")
	cmd.Flags().Int("timeout", 2000, "timeout in milliseconds")
	appendCommand(cmd)
}

func runTcpScan(cmd *cobra.Command, _ []string) {
	target, _ := cmd.Flags().GetString("ip")
	port1, _ := cmd.Flags().GetInt("port1")
	port2, _ := cmd.Flags().GetInt("port2")
	timeout, _ := cmd.Flags().GetInt("timeout") // 毫秒超时

	ports := make([]int, 0, port2-port1+1)
	for p := port1; p <= port2; p++ {
		ports = append(ports, p)
	}
	total := len(ports)

	// 收集开放端口
	var openPorts []int
	var mu sync.Mutex

	var wg sync.WaitGroup
	var scanned int64
	maxConcurrency := 1000

	semaphore := make(chan struct{}, maxConcurrency)

	// 进度显示使用 stderr
	go func() {
		for {
			done := atomic.LoadInt64(&scanned)
			percent := float64(done) / float64(total) * 100
			fmt.Fprintf(os.Stderr, "\r[%d/%d] %.1f%%", done, total, percent)
			if done >= int64(total) {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	start := time.Now()
	for _, port := range ports {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(p int) {
			defer func() {
				<-semaphore
				wg.Done()
				atomic.AddInt64(&scanned, 1)
			}()
			addr := fmt.Sprintf("%s:%d", target, p)
			conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Millisecond)
			if err != nil {
				return
			}
			defer conn.Close()
			// 结果输出到 stdout
			fmt.Printf("%s:%d is open\n", target, p)
			// 收集到汇总
			mu.Lock()
			openPorts = append(openPorts, p)
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Fprintf(os.Stderr, "\r[%d/%d] 100.0%%\n", total, total)

	// 打印汇总
	fmt.Printf("\n=== 扫描汇总 ===\n")
	fmt.Printf("目标: %s\n", target)
	fmt.Printf("端口范围: %d - %d\n", port1, port2)
	fmt.Printf("开放端口数: %d/%d\n", len(openPorts), total)
	if len(openPorts) > 0 {
		fmt.Printf("开放端口: %v\n", openPorts)
	}
	fmt.Printf("耗时: %v\n", elapsed)
}

func init() {
	cmd := &cobra.Command{
		Use:   "net:scan-range",
		Short: "TCP IP段端口扫描工具",
		Run:   runTcpScan2,
		// Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	}
	cmd.Flags().String("ip1", "127.0.0.1", "IP address 1")
	cmd.Flags().String("ip2", "127.0.0.1", "IP address 2")
	cmd.Flags().String("port", "22,80", "port number 22,80")
	appendCommand(cmd)
}

func runTcpScan2(cmd *cobra.Command, _ []string) {
	start := time.Now()
	startIP, _ := cmd.Flags().GetString("ip1")
	endIP, _ := cmd.Flags().GetString("ip2")
	portListStr, _ := cmd.Flags().GetString("port")

	port := lo.Map(strings.Split(portListStr, ","), func(t string, _ int) int {
		return cast.ToInt(t)
	})

	scanIPRange(startIP, endIP, port)

	elapsed := time.Since(start)
	fmt.Printf("Scan completed in %v\n", elapsed)
}

func scanIPRange(startIP, endIP string, ports []int) {
	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)

	if start.To4() == nil || end.To4() == nil {
		fmt.Println("Invalid IP address")
		return
	}

	startInt := ipToInt(start.To4())
	endInt := ipToInt(end.To4())

	if startInt > endInt {
		fmt.Println("Invalid IP range")
		return
	}

	totalIPs := int(endInt - startInt + 1)
	total := totalIPs * len(ports)

	var wg sync.WaitGroup
	var scanned int64
	maxConcurrency := 1000

	semaphore := make(chan struct{}, maxConcurrency)

	// 进度显示
	go func() {
		for {
			done := atomic.LoadInt64(&scanned)
			percent := float64(done) / float64(total) * 100
			fmt.Fprintf(os.Stderr, "\r[%d/%d] %.1f%%", done, total, percent)
			if done >= int64(total) {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 遍历 IP 段中的所有 IP 地址，依次进行端口扫描
	for i := startInt; i <= endInt; i++ {
		ip := intToIP(i)

		wg.Add(len(ports))

		// 对指定 IP 地址上的所有端口依次进行扫描
		for _, port := range ports {
			semaphore <- struct{}{}
			go func(ip string, port int) {
				defer func() {
					<-semaphore
					wg.Done()
					atomic.AddInt64(&scanned, 1)
				}()
				target := fmt.Sprintf("%s:%d", ip, port)
				conn, err := net.DialTimeout("tcp", target, 200*time.Millisecond)
				if err != nil {
					return
				}
				defer conn.Close()
				fmt.Printf("%s:%d is open\n", ip, port)
			}(ip, port)
		}
	}

	wg.Wait()
	fmt.Fprintf(os.Stderr, "\r[%d/%d] 100.0%%\n", total, total)
	fmt.Println("Scan completed")
}

func ipToInt(ip net.IP) int64 {
	if len(ip) == 16 {
		return int64(ip[12])<<24 | int64(ip[13])<<16 | int64(ip[14])<<8 | int64(ip[15])
	}
	return int64(ip[0])<<24 | int64(ip[1])<<16 | int64(ip[2])<<8 | int64(ip[3])
}

func intToIP(i int64) string {
	ip := make(net.IP, 4)
	ip[0] = byte(i >> 24 & 0xFF)
	ip[1] = byte(i >> 16 & 0xFF)
	ip[2] = byte(i >> 8 & 0xFF)
	ip[3] = byte(i & 0xFF)
	return ip.String()
}
