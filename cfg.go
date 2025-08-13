package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	AppName        = "<unknown app name>"
	Version        = "<unknown build version>"
	BuildTimeStamp = "<unknown build timestamp>"
	versionFlag    *bool
)

func Configure() {
	versionFlag = flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		ShowVersion()
	}
}

func ShowVersion() {

	// Get the name of the application from the file path
	_, file, _, _ := runtime.Caller(2)
	AppName = strings.Split(file, "/")[len(strings.Split(file, "/"))-2]

	if Version == "" {
		Version = "<unknown git version>"
	}
	btime, err := time.Parse("2006-01-02_15:04:05_MST", BuildTimeStamp)
	if err != nil {
		btime = time.Now()
	}
	fmt.Printf("%s - version %s built on %s \n", AppName, Version, btime.Format("02.01.2006 - 15:04:05 MST"))
	os.Exit(0)
}
