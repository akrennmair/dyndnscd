package main

import (
	"goconf.googlecode.com/hg"
	"flag"
	"os"
	"fmt"
	"time"
	"runtime"
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

	logchan := make(chan LogMsg)
	go Logger(c, logchan)

	SpawnPollers(c, logchan)

	go func() {
		for {
			time.Sleep(int64(120000000000))
			logchan <- NewLogMsg(INFO, "memory: " + MemoryStatistics())
		}
	}()

	done := make(chan int)
	<-done;
}

func MemoryStatistics() string {
	var p []runtime.MemProfileRecord
	n, ok := runtime.MemProfile(nil, false)
	for {
		p = make([]runtime.MemProfileRecord, n+50)
		n, ok = runtime.MemProfile(p, false)
		if ok {
			p = p[0:n]
			break
		}
	}

	var total runtime.MemProfileRecord
	for i := range p {
		r := &p[i]
		total.AllocBytes += r.AllocBytes
		total.AllocObjects += r.AllocObjects
		total.FreeBytes += r.FreeBytes
		total.FreeObjects += r.FreeObjects
	}

	return fmt.Sprintf("%d in use objects (%d in use bytes) | Alloc: %d TotalAlloc: %d",
		total.InUseObjects(), total.InUseBytes(), runtime.MemStats.Alloc, runtime.MemStats.TotalAlloc)
}
