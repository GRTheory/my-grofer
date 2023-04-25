package general

import (
	"context"
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

// ServeDiskRates serves the disk rate data to the data channel.
func ServeDiskRates(ctx context.Context, dataChannel chan AggregateMetrics) error {
	var partitions []disk.PartitionStat
	var err error
	// in this situation, we choose to use separased disk information.
	partitions, err = disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return err
	}
	rows := [][]string{{"Mount", "Total", "Used %", "Used", "Free", "FS Type"}}
	for _, value := range partitions {
		usageVals, _ := disk.UsageWithContext(ctx, value.Mountpoint)

		if strings.HasPrefix(value.Device, "/dev/loop") {
			continue
		} else if strings.HasPrefix(value.Mountpoint, "/var/lib/docker") {
			continue
		} else {
			path := usageVals.Path
			total := fmt.Sprintf("%.2f G", float64(usageVals.Total)/(1024*1024*1024))
			used := fmt.Sprintf("%.2f G", float64(usageVals.Used)/(1024*1024*1024))
			usedPercent := fmt.Sprintf("%.2f %s", usageVals.UsedPercent, "%")
			free := fmt.Sprintf("%.2f G", float64(usageVals.Free)/(1024*1024*1024))
			fs := usageVals.Fstype
			row := []string{path, total, usedPercent, used, free, fs}
			rows = append(rows, row)
		}
	}

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

// ServeNetRates serves info about the network to the data channel.
func ServeNetRates(ctx context.Context, dataChannel chan AggregateMetrics) error {
	netStats, err := net.IOCountersWithContext(ctx, false)
	if err != nil {
		return err
	}
	IO := make(map[string][]float64)
	for _, IOStat := range netStats {
		nic := []float64{float64(IOStat.BytesSent), float64(IOStat.BytesRecv)}
		IO[IOStat.Name] = nic
	}

	data := AggregateMetrics{
		FieldSet: "NET",
		NetStats: IO,
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case dataChannel <- data:
		return nil
	}
}
