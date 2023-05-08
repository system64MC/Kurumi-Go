// go: generate goversioninfo -icon = kuruicon.ico
package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/system64MC/Kurumi-Go/kuruApp"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	kuruApp.Init()

	// openbrowser("https://www.youtube.com/watch?v=xvFZjo5PgG0")
}
