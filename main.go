package main

import (
	"goconf.googlecode.com/hg"
	"flag"
	"os"
	"fmt"
)

func main() {
	var configfile *string = flag.String("f", "", "configuration file")
	flag.Parse()

	if *configfile == "" {
		fmt.Println("usage: dyndnscd -f <configfile>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	c, err := conf.ReadConfigFile(*configfile)
	if err != nil {
		fmt.Printf("reading configuration file failed: %v\n", err)
		os.Exit(1)
	}

	SpawnPollers(c)

	done := make(chan int)
	<-done;
}
