package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type updater interface {
	UpdateIP(ctx context.Context, ip net.IP) error
	Target() string
}

type urlUpdater struct {
	url string
}

func newUpdater(conf configSection) (updater, error) {
	if conf.UpdateURL == "" {
		return nil, fmt.Errorf("%s: missing update_url", conf.Name)
	}
	return &urlUpdater{
		url: conf.UpdateURL,
	}, nil
}

func (u *urlUpdater) Target() string {
	return u.url
}

func (u *urlUpdater) UpdateIP(ctx context.Context, ip net.IP) error {
	fullURL := strings.Replace(u.url, "<ip>", ip.String(), -1)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
