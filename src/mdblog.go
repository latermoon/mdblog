package main

import (
	"flag"
	"log"
	"runtime"
	"server"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// go run src/mdblog.go -h localhost:8801 -d /home/workspace/blog.latermoon.me
func main() {
	host, dir := readFlags()
	if dir == "" {
		log.Fatal("no -d website")
	}

	log.Printf("%s mdblog sercive start\n", host)
	server.ListenAndServe(host, dir)
}

func readFlags() (host string, dir string) {
	flag.StringVar(&dir, "d", "", "website dir")
	flag.StringVar(&host, "h", "localhost:8801", "mdblog background service")
	flag.Parse()
	return
}
