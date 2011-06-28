package main

import (
	"fmt"
	"time"
	"goconf.googlecode.com/hg"
	"http"
	"bufio"
	"strings"
	"syscall"
	"net"
	"os"
	"unsafe"
)

func SpawnPollers(c *conf.ConfigFile) {
	sections := c.GetSections()
	for i:= range sections {
		sectionname := sections[i]
		sectiontype, _ := c.GetString(sectionname, "type")
		if sectiontype == "" {
			continue
		}
		switch sectiontype {
			case "ipbouncer":
				bouncer_url, _ := c.GetString(sectionname, "bouncer_url")
				update_url, _ := c.GetString(sectionname, "update_url")
				interval, _ := c.GetInt(sectionname, "interval")
				go IpBouncerPoller(bouncer_url, update_url, interval)
			case "device":
				device, _ := c.GetString(sectionname, "device")
				update_url, _ := c.GetString(sectionname, "update_url")
				interval, _ := c.GetInt(sectionname, "interval")
				go DevicePoller(device, update_url, interval)
			default:
				fmt.Printf("Warning: unknown type %v\n", sectiontype)
		}
	}
}

func IpBouncerPoller(bouncer_url string, update_url string, interval int) {
	for {
		fmt.Printf("polling %s\n", bouncer_url)
		httpclient := new(http.Client)
		resp, err := httpclient.Get(bouncer_url)
		if err == nil {
			httpdata := bufio.NewReader(resp.Body)
			ipdata, _, err2 := httpdata.ReadLine()
			if err2 == nil {
				ip := string(ipdata)
				UpdateIP(update_url, ip)
			}
			resp.Body.Close()
		}
		time.Sleep(int64(interval) * int64(1000000000))
	}
}

func DevicePoller(device string, update_url string, interval int) {
	for {
		fmt.Printf("polling %s\n", device)

		ip, err := GetIPFromDevice(device)
		if err == nil {
			UpdateIP(update_url, ip.String())
		} else {
			fmt.Printf("device poller error: %v\n", err)
		}

		time.Sleep(int64(interval) * int64(1000000000))
	}
}

func UpdateIP(update_url string, ip string) {
	full_url := strings.Replace(update_url, "<ip>", ip, -1)
	fmt.Printf("full_url = %v ip = %v\n", full_url, ip)
	httpclient := new(http.Client)
	resp, err := httpclient.Get(full_url)
	if err != nil {
		fmt.Printf("an error occured while updating IP: %v\n", err)
	}
	resp.Body.Close()
}

func GetIPFromDevice(device string) (ip net.IP, err os.Error) {
	var ifreqbuf [40]byte

	for i := 0 ; i < 40 ; i++ {
		ifreqbuf[i] = 0
	}

	for i := 0 ; i < len(device) ; i++ {
		ifreqbuf[i] = device[i]
	}

	socketfd, _, _ := syscall.Syscall(syscall.SYS_SOCKET, syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	defer syscall.Close(int(socketfd))

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(socketfd), uintptr(syscall.SIOCGIFADDR), uintptr(unsafe.Pointer(&ifreqbuf)))
	if err = os.NewSyscallError("SYS_IOCTL", int(errno)); err != nil {
		return
	}
	ip = net.IPv4(ifreqbuf[20], ifreqbuf[21], ifreqbuf[22], ifreqbuf[23])
	return
}
