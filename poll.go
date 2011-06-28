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
	for i := range sections {
		sectionname := sections[i]
		if sectionname == "default" {
			continue
		}

		fetcher, err1 := NewFetcher(c, sectionname)
		if err1 != nil {
			fmt.Printf("couldn't create fetcher for section %s: %v\n", sectionname, err1)
			continue
		}

		updater, err2 := NewUpdater(c, sectionname)
		if err2 != nil {
			fmt.Printf("couldn't create updater for section %s: %v\n", sectionname, err2)
			continue
		}

		interval, err3 := c.GetInt(sectionname, "interval")
		if err3 != nil {
			fmt.Printf("couldn't get interval for section %s: %v\n", sectionname, err3)
			continue
		}

		go Poller(fetcher, updater, interval)
	}
}

func Poller(f Fetcher, u Updater, interval int) {
	old_ip := net.IPv4(0, 0, 0, 0)
	for {
		ip, err := f.FetchIP()
		if err != nil {
			fmt.Printf("Fetching IP failed: %v\n", err)
		} else if (!ip.Equal(old_ip)) {
			if err2 := u.UpdateIP(ip); err != nil {
				fmt.Printf("Updating IP failed: %v\n", err2)
			} else {
				old_ip = ip
			}
		}
		time.Sleep(int64(interval) * int64(1000000000))
	}
}
