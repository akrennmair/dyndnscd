package main

import (
	"goconf.googlecode.com/hg"
	"os"
	"net"
	"http"
	"strings"
)

type Updater interface {
	UpdateIP(ip net.IP) os.Error
	Target() string
}

type URLUpdater struct {
	url string
}

func NewUpdater(c *conf.ConfigFile, section string) (Updater, os.Error) {
	u := &URLUpdater{}
	update_url, err := c.GetString(section, "update_url")
	if err != nil {
		return nil, err
	}
	u.url = update_url
	return u, nil
}

func (u URLUpdater) Target() string {
		return u.url
}

func (u URLUpdater) UpdateIP(ip net.IP) os.Error {
	full_url := strings.Replace(u.url, "<ip>", ip.String(), -1)
	httpclient := &http.Client{}
	resp, err := httpclient.Get(full_url)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
