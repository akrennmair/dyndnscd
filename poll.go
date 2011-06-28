package main

import (
	"fmt"
	"time"
	"goconf.googlecode.com/hg"
	"http"
	"bufio"
	"strings"
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
