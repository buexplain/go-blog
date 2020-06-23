package c_sysInfo

import (
	"fmt"
	a_boot "github.com/buexplain/go-blog/app/boot"
	"github.com/buexplain/go-blog/helpers"
	"github.com/buexplain/go-fool"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

func Index(ctx *fool.Ctx, w *fool.Response, r *fool.Request) error {
	result := make(map[string]interface{})

	//获取基础信息
	base := map[string]string{}
	if v, err := host.Info(); err == nil {
		base["hostName"] = v.Hostname
		base["os"] = v.OS
	} else {
		return fmt.Errorf("读取系统信息失败: %w", err)
	}
	base["goVersion"] = runtime.Version()
	base["pid"] = strconv.Itoa(os.Getpid())
	base["time"] = time.Now().Format(time.RFC3339)
	result["base"] = base

	//读取内存信息
	if v, err := mem.VirtualMemory(); err == nil {
		tmp := map[string]string{}
		tmp["total"] = helpers.FormatSize(int64(v.Total))
		tmp["free"] = helpers.FormatSize(int64(v.Available))
		tmp["used"] = helpers.FormatSize(int64(v.Used))
		result["memory"] = tmp
	} else {
		return fmt.Errorf("读取内存信息失败: %w", err)
	}

	//读取磁盘信息
	if v, err := disk.Usage(a_boot.ROOT_PATH); err == nil {
		tmp := map[string]string{}
		tmp["total"] = helpers.FormatSize(int64(v.Total))
		tmp["free"] = helpers.FormatSize(int64(v.Free))
		tmp["used"] = helpers.FormatSize(int64(v.Used))
		result["disk"] = tmp
	} else {
		return fmt.Errorf("读取磁盘信息失败: %w", err)
	}

	//读取cpu信息
	if v, err := cpu.Info(); err == nil {
		tmp := make([]struct {
			Name  string
			Cores int32
		}, 0)
		for _, info := range v {
			tmp = append(tmp, struct {
				Name  string
				Cores int32
			}{Name: info.ModelName, Cores: info.Cores})
		}
		result["cpu"] = tmp
	} else {
		return fmt.Errorf("读取cpu信息失败: %w", err)
	}

	//网卡
	if v, err := net.IOCounters(true); err == nil {
		tmp := make([]struct {
			Name string
			Recv string
			Sent string
		}, 0)
		for _, info := range v {
			tmp = append(tmp, struct {
				Name string
				Recv string
				Sent string
			}{Name: info.Name, Recv: helpers.FormatSize(int64(info.BytesRecv)), Sent: helpers.FormatSize(int64(info.BytesSent))})
		}
		result["io"] = tmp
	} else {
		return fmt.Errorf("读取网卡信息失败: %w", err)
	}

	//读取负载信息
	if v, err := load.Avg(); err == nil {
		tmp := map[string]interface{}{}
		tmp["one"] = v.Load1
		tmp["five"] = v.Load5
		tmp["fifteen"] = v.Load15
		result["load"] = tmp
	} else {
		if runtime.GOOS == "windows" {
			tmp := map[string]interface{}{}
			tmp["one"] = 0
			tmp["five"] = 0
			tmp["fifteen"] = 0
			result["load"] = tmp
		} else {
			return fmt.Errorf("读取负载信息失败: %w", err)
		}
	}

	return w.Assign("result", result).View(http.StatusOK, "backend/server/sysInfo/index.html")
}
