package main

import (
	"goconf.googlecode.com/hg"
	"os"
	"net"
	"bufio"
	"http"
	"syscall"
	"unsafe"
)

type Fetcher interface {
	FetchIP() (net.IP, os.Error)
	Source() string
}

func NewFetcher(c *conf.ConfigFile, section string) (f Fetcher, e os.Error) {
	sectiontype, err := c.GetString(section, "type")
	if err != nil {
		return nil, err
	}
	switch sectiontype {
		case "ipbouncer":
			bouncer_url, e1 := c.GetString(section, "bouncer_url")
			if e1 != nil {
				return nil, NewConfigMissingError("bouncer_url")
			}
			return NewIPBouncerFetcher(bouncer_url), nil
		case "device":
			device, e1 := c.GetString(section, "device")
			if e1 != nil {
				return nil, NewConfigMissingError("device")
			}
			return NewDeviceFetcher(device), nil
	}
	return nil, NewUnknownSectionTypeError(sectiontype)
}

type IPBouncerFetcher struct {
	bouncer_url string
}

func NewIPBouncerFetcher(bouncer_url string) Fetcher {
	f := new(IPBouncerFetcher)
	f.bouncer_url = bouncer_url
	return f
}

func (f IPBouncerFetcher) Source() string {
	return f.bouncer_url
}

func (f IPBouncerFetcher) FetchIP() (ip net.IP, e os.Error) {
	httpclient := new(http.Client)
	resp, err := httpclient.Get(f.bouncer_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	httpdata := bufio.NewReader(resp.Body)
	ipdata, _, err2 := httpdata.ReadLine()
	if err2 != nil {
		return nil, err2
	}

	return net.ParseIP(string(ipdata)), nil
}

type DeviceFetcher struct {
	device string
}

func NewDeviceFetcher(device string) Fetcher {
	f := new(DeviceFetcher)
	f.device = device
	return f
}

func (f DeviceFetcher) Source() string {
	return f.device
}

func (f DeviceFetcher) FetchIP() (ip net.IP, err os.Error) {
	var ifreqbuf [40]byte

	for i := 0 ; i < 40 ; i++ {
		ifreqbuf[i] = 0
	}

	for i := 0 ; i < len(f.device) ; i++ {
		ifreqbuf[i] = f.device[i]
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
