package main

import (
	"flag"
	"ledis"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var configFile = flag.String("config", "/etc/ledis.json", "ledisdb config file")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if len(*configFile) == 0 {
		println("must use a config file")
		return
	}

	cfg, err := server.NewConfigWithFile(*configFile)
	if err != nil {
		println(err.Error())
		return
	}

	var app *server.App
	app, err = server.NewApp(cfg)
	if err != nil {
		println(err.Error())
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-sc

		app.Close()
	}()

	app.Run()
}
