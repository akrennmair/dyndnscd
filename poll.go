package main

import (
	"fmt"
	"github.com/akrennmair/goconf"
	"net"
	"time"
)

func SpawnPollers(c *conf.ConfigFile, logchan chan LogMsg) {
	sections := c.GetSections()
	for i := range sections {
		sectionname := sections[i]
		if sectionname == "default" {
			continue
		}

		fetcher, err1 := NewFetcher(c, sectionname)
		if err1 != nil {
			logchan <- NewLogMsg(WARN, fmt.Sprintf("couldn't create fetcher for section %s: %v", sectionname, err1))
			continue
		}

		updater, err2 := NewUpdater(c, sectionname)
		if err2 != nil {
			logchan <- NewLogMsg(WARN, fmt.Sprintf("couldn't create updater for section %s: %v", sectionname, err2))
			continue
		}

		interval, err3 := c.GetInt(sectionname, "interval")
		if err3 != nil {
			logchan <- NewLogMsg(WARN, fmt.Sprintf("couldn't get interval for section %s: %v", sectionname, err3))
			continue
		}

		go Poller(sectionname, fetcher, updater, interval, logchan)
	}
}

func Poller(section string, f Fetcher, u Updater, interval int, logchan chan LogMsg) {
	logchan <- NewLogMsg(DEBUG, "started Poller for section "+section)
	old_ip := net.IPv4(0, 0, 0, 0)
	for {
		ip, err := f.FetchIP()
		if err != nil {
			logchan <- NewLogMsg(ERROR, fmt.Sprintf("%s: fetching IP from %s failed: %v", section, f.Source(), err))
		} else {
			logchan <- NewLogMsg(DEBUG, fmt.Sprintf("%s: fetched IP %v", section, ip))
			if !ip.Equal(old_ip) {
				if err2 := u.UpdateIP(ip); err != nil {
					logchan <- NewLogMsg(ERROR, fmt.Sprintf("Updating IP to %s failed: %v", u.Target(), err2))
				} else {
					old_ip = ip
				}
			} else {
				logchan <- NewLogMsg(DEBUG, section+": new IP same as old IP")
			}
		}
		logchan <- NewLogMsg(DEBUG, fmt.Sprintf("%s: sleeping for %d seconds", section, interval))
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
