package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/fraugster/cli"
	"gopkg.in/yaml.v2"
)

func main() {
	ctx := cli.Context()

	dialer := &net.Dialer{}

	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, _, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp4", addr)
	}

	configfile := flag.String("f", "", "configuration file")
	flag.Parse()

	if *configfile == "" {
		fmt.Println("usage: dyndnscd -f <configfile>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	f, err := os.Open(*configfile)
	if err != nil {
		log.Fatalf("Opening configuration file %s failed: %v", *configfile, err)
	}
	defer f.Close()

	var conf config

	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		log.Fatalf("Parsing configuration file %s failed: %v", *configfile, err)
	}

	spawnPollers(ctx, conf)

	<-ctx.Done()
}

type config []configSection

type configSection struct {
	Name       string        `yaml:"name"`
	Interval   time.Duration `yaml:"duration"`
	Type       string        `yaml:"type"`
	BouncerURL string        `yaml:"bouncer_url"`
	Device     string        `yaml:"device"`
	UpdateURL  string        `yaml:"update_url"`
}
