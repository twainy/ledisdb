package server

import (
	ledis_client "ledis/client"
	"os"
	"sync"
	"testing"
)

var testAppOnce sync.Once
var testApp *App

var testLedisClient *ledis_client.Client

func newTestLedisClient() {
	cfg := new(ledis_client.Config)
	cfg.Addr = "127.0.0.1:16380"
	cfg.MaxIdleConns = 4
	testLedisClient = ledis_client.NewClient(cfg)
}

func getTestConn() *ledis_client.Conn {
	startTestApp()
	return testLedisClient.Get()
}

func startTestApp() {
	f := func() {
		newTestLedisClient()

		os.RemoveAll("/tmp/testdb")

		var d = []byte(`
            {
                "data_dir" : "/tmp/testdb",
                "addr" : "127.0.0.1:16380",
                "db" : {        
                    "compression":true,
                    "block_size" : 32768,
                    "write_buffer_size" : 2097152,
                    "cache_size" : 20971520,
                    "max_open_files" : 1024
                }    
            }
            `)

		cfg, err := NewConfig(d)
		if err != nil {
			println(err.Error())
			panic(err)
		}

		testApp, err = NewApp(cfg)
		if err != nil {
			println(err.Error())
			panic(err)
		}

		go testApp.Run()
	}

	testAppOnce.Do(f)
}

func TestApp(t *testing.T) {
	startTestApp()
}
