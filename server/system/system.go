package system

import (
	"fmt"
	"net"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	psutilNet "github.com/shirou/gopsutil/v3/net"
)

type SystemInformation struct {
	VirtualMem *mem.VirtualMemoryStat        `json:"virtual_memory"`
	SwapMem    *mem.SwapMemoryStat           `json:"swap_memory"`
	SwapDev    []*mem.SwapDevice             `json:"swap_devices"`
	CPUInfo    []cpu.InfoStat                `json:"cpu_info"`
	NetInfo    []psutilNet.ProtoCountersStat `json:"net_info"`
	HostInfo   *host.InfoStat                `json:"host_info"`
}

// Relevant for engineer routes.
// Right now we will only fetch Linux information since that is the recommended host OS.
func FetchSystemInformation() (*SystemInformation, error) {
	if strings.ToLower(runtime.GOOS) != "linux" {
		return nil, fmt.Errorf("running %s, ignoring system information", runtime.GOOS)
	}

	virtualMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	swapMem, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	swapDevices, err := mem.SwapDevices()
	if err != nil {
		return nil, err
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	netStats, err := psutilNet.ProtoCounters([]string{"tcp", "udp"})
	if err != nil {
		return nil, err
	}

	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	var systemInfo SystemInformation

	systemInfo.VirtualMem = virtualMem
	systemInfo.SwapMem = swapMem
	systemInfo.SwapDev = swapDevices
	systemInfo.CPUInfo = cpuInfo
	systemInfo.NetInfo = netStats
	systemInfo.HostInfo = hostInfo

	return &systemInfo, nil
}

func GetIPForDynamicHost(domain string) (net.IP, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	return ips[0], nil
}
