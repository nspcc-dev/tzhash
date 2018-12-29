package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/nspcc-dev/tzhash/tz"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	filename   = flag.String("name", "-", "file to use")
)

func main() {
	var (
		f   io.Reader
		err error
		h   = tz.New()
	)

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *filename != "-" {
		if f, err = os.Open(*filename); err != nil {
			log.Fatal("could not open file: ", err)
		}
	} else {
		f = os.Stdin
	}

	if _, err := io.Copy(h, f); err != nil {
		log.Fatal("error while reading file: ", err)
	}
	fmt.Printf("%x\t%s\n", h.Sum(nil), *filename)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
