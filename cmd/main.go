package main

import (
	"os"
	"github.com/adrianduke/crawler"
)

func main() {
	os.Exit(crawler.EntryPoint(os.Args[1:], os.Stdout, os.Stderr))
}
