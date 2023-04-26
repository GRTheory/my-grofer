package general

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type MonitorInterface interface {
	String() string
}

func roundOff(num uint64) float64 {
	x := float64(num) / (1024 * 1024 * 1024)
	return math.Round(x*10) / 10
}

// ServeInfo provides information about the system such as OS info, uptime, boot time, etc.
func ServeInfo(ctx context.Context, dataChannel chan AggregateMetrics) error {
	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return err
	}

	hostInfo := [][]string{
		{"Hostname", info.Hostname},
		{"Processes", fmt.Sprintf("%d", info.Procs)},
		{"OS/Platform", fmt.Sprintf("%s/%s %s", info.OS, info.Platform, info.PlatformVersion)},
		{"Kernel/Arch", fmt.Sprintf("%s/%s", info.KernelVersion, info.KernelArch)},
	}

	data := AggregateMetrics{
		FieldSet: "INFO",
		HostInfo: hostInfo,
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case dataChannel <- data:
		return nil
	}
}

// GetCPURates fetches and returns teh current cpu rate
func GetCPURates(ctx context.Context) ([]float64, error) {
	cpuRates, err := cpu.PercentWithContext(ctx, time.Second, true)
	if err != nil {
		return nil, err
	}
	return cpuRates, nil
}

// ServeCPURates serves the cpu rates to the cpu channel
func ServeCPURates(ctx context.Context, dataChannel chan AggregateMetrics) error {
	cpuRates, err := cpu.PercentWithContext(ctx, time.Second, true)
	if err != nil {
		return err
	}
	data := AggregateMetrics{
		FieldSet: "CPU",
		CPUStats: cpuRates,
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case dataChannel <- data:
		return nil
	}
}

// ServeMemRates serves stats about the memory to the data channel.
func ServeMemRates(ctx context.Context, dataChannel chan AggregateMetrics) error {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	memRates := []float64{roundOff(memory.Total), roundOff(memory.Used), roundOff(memory.Available), roundOff(memory.Free), roundOff(memory.Cached)}

	data := AggregateMetrics{
		FieldSet: "MEM",
		MemStats: memRates,
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case dataChannel <- data:
		return nil
	}
}

type DiskMonitor struct {
	MountPath   string
	Total       float64
	UsedPercent float64
	Used        float64
	Free        float64
	FsType      string
}

// ServeDiskRates serves the disk rate data to the data channel.
func ServeDiskRates(ctx context.Context, dataChannel chan AggregateMetrics) error {
	var partitions []disk.PartitionStat
	var err error
	// in this situation, we choose to use separased disk information.
	partitions, err = disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return err
	}
	// rows := [][]string{{"Mount", "Total", "Used %", "Used", "Free", "FS Type"}}
	rows := make([]DiskMonitor, 0, len(partitions))
	for _, value := range partitions {
		usageVals, _ := disk.UsageWithContext(ctx, value.Mountpoint)

		if strings.HasPrefix(value.Device, "/dev/loop") {
			continue
		} else if strings.HasPrefix(value.Mountpoint, "/var/lib/docker") {
			continue
		} else {
			tempDiskMonitor := DiskMonitor{
				MountPath:   usageVals.Path,
				Total:       float64(usageVals.Total) / (1024 * 1024 * 1024),
				UsedPercent: usageVals.UsedPercent,
				Used:        float64(usageVals.Used) / (1024 * 1024 * 1024),
				Free:        float64(usageVals.Free) / (1024 * 1024 * 1024),
				FsType:      usageVals.Fstype,
			}
			rows = append(rows, tempDiskMonitor)
		}
	}

	result, err := json.Marshal(rows)
	if err != nil {
		return err
	}

	fmt.Println(string(result))

	data := AggregateMetrics{
		FieldSet:  "DISK",
		DiskStats: rows,
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case dataChannel <- data:
		return nil
	}
}

type NetMonitor struct {
	Name          string
	SendBytes     float64
	ReceivedBytes float64
}

// ServeNetRates gathers network related metrics and sends them to the dataChannel. 
// It takes the context and  the dataChannel as arguments and returns an error if 
// one encountered. The function gathers information such as bytes sent and received 
// and stores them in a NetMonitor struct which is then sent to the dataChannel. 
// If a context is done, an error is returned.
func ServeNetRates(ctx context.Context, dataChannel chan AggregateMetrics) error {
	netStats, err := net.IOCountersWithContext(ctx, false)
	if err != nil {
		return err
	}
	netMonitors := make([]NetMonitor, 0, len(netStats))
	for _, ioStat := range netStats {
		tempNetMonitor := NetMonitor{
			Name:          ioStat.Name,
			SendBytes:     float64(ioStat.BytesSent),
			ReceivedBytes: float64(ioStat.BytesRecv),
		}
		netMonitors = append(netMonitors, tempNetMonitor)
	}
	data := AggregateMetrics{
		FieldSet: "NET",
		NetStats: netMonitors,
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case dataChannel <- data:
		return nil
	}
}
