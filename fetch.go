package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"syscall"
	"unsafe"
)

type fetcher interface {
	FetchIP(context.Context) (net.IP, error)
	Source() string
}

func newFetcher(conf configSection) (fetcher, error) {
	switch conf.Type {
	case "ipbouncer":
		if conf.BouncerURL == "" {
			return nil, fmt.Errorf("%s: missing bouncer_url", conf.Name)
		}
		return newIPBouncerFetcher(conf.BouncerURL), nil
	case "device":
		if conf.Device == "" {
			return nil, fmt.Errorf("%s: missing device", conf.Name)
		}
		return newDeviceFetcher(conf.Device), nil
	default:
		return nil, fmt.Errorf("unknown type %s", conf.Type)
	}
}

type ipBouncerFetcher struct {
	url string
}

func newIPBouncerFetcher(url string) fetcher {
	return &ipBouncerFetcher{
		url: url,
	}
}

func (f *ipBouncerFetcher) Source() string {
	return f.url
}

func (f *ipBouncerFetcher) FetchIP(ctx context.Context) (net.IP, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", f.url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned HTTP %d", f.url, resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(data))
	if ip == nil {
		return nil, fmt.Errorf("parsing IP %s failed", string(data))
	}

	return ip, nil
}

type deviceFetcher struct {
	device string
}

func newDeviceFetcher(device string) fetcher {
	return &deviceFetcher{
		device: device,
	}
}

func (f *deviceFetcher) Source() string {
	return f.device
}

func (f *deviceFetcher) FetchIP(context.Context) (net.IP, error) {
	var ifreqbuf [40]byte

	for i := 0; i < len(f.device); i++ {
		ifreqbuf[i] = f.device[i]
	}

	socketfd, _, errno := syscall.Syscall(syscall.SYS_SOCKET, syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err := os.NewSyscallError("SYS_SOCKET", errno); err != nil {
		return nil, err
	}
	defer syscall.Close(int(socketfd))

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, uintptr(socketfd), uintptr(syscall.SIOCGIFADDR), uintptr(unsafe.Pointer(&ifreqbuf)))
	if err := os.NewSyscallError("SYS_IOCTL", errno); err != nil {
		return nil, err
	}
	return net.IPv4(ifreqbuf[20], ifreqbuf[21], ifreqbuf[22], ifreqbuf[23]), nil
}
