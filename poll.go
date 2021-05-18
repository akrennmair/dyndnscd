package main

import (
	"context"
	"log"
	"net"
	"time"
)

func spawnPollers(ctx context.Context, conf config) {
	for _, sectionConf := range conf {

		fetcher, err := newFetcher(sectionConf)
		if err != nil {
			log.Printf("Couldn't create fetcher for configuration %s: %v", sectionConf.Name, err)
			continue
		}

		updater, err := newUpdater(sectionConf)
		if err != nil {
			log.Printf("Couldn't create updater for configuration %s: %v", sectionConf.Name, err)
			continue
		}

		go poller(ctx, sectionConf, fetcher, updater)
	}
}

func poller(ctx context.Context, conf configSection, f fetcher, u updater) {
	log.Printf("Started poller for section %s", conf.Name)
	oldIP := net.IPv4(0, 0, 0, 0)

	interval := conf.Interval
	if interval == 0 {
		interval = 60 * time.Second
	}

	ticker := time.NewTicker(interval)

	poll := func() {
		ip, err := f.FetchIP(ctx)
		if err != nil {
			log.Printf("%s: fetching IP from %s failed: %v", conf.Name, f.Source(), err)
			return
		}
		if ip.Equal(oldIP) {
			log.Printf("%s: new IP same as old IP", conf.Name)
			return
		}
		if err := u.UpdateIP(ctx, ip); err != nil {
			log.Printf("%s: updating IP to %s failed: %v", conf.Name, u.Target(), err)
			return
		}
		oldIP = ip
	}

	poll()

	for {
		select {
		case <-ticker.C:
			poll()
		case <-ctx.Done():
			return
		}
	}
}
