package main

import (
	"context"
	"fmt"
	"os"

	"github.com/QOSGroup/cassini/adapter/ports"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/event"
	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	"github.com/spf13/viper"
)

// 命令行 events 命令执行方法
var events = func() (cancel context.CancelFunc, err error) {
	conf := config.GetConfig()
	var cancels []context.CancelFunc
	var cancelFunc context.CancelFunc

	errChannel := make(chan error, 1)
	startLog(errChannel)
	startExporter(errChannel)

	for _, mockConf := range conf.Mocks {
		cancelFunc, err = subscribe(mockConf.RPC.NodeAddress, mockConf.Subscribe)
		if err != nil {
			return
		}
		cancels = append(cancels, cancelFunc)
	}
	cancel = func() {
		for _, cancelJob := range cancels {
			if cancelJob != nil {
				cancelJob()
			}
		}
		log.Debug("Cancel events subscribe service")
	}
	return
}

//subscribe 从websocket服务端订阅event
//remote 服务端地址 example  "127.0.0.1:27657"
func subscribe(remote string, query string) (context.CancelFunc, error) {
	ip, port, err := ports.ParseNodeAddress(remote)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conf := &ports.AdapterConfig{
		ChainName: "cassini-events",
		ChainType: "qos",
		IP:        ip,
		Port:      port,
		Query:     query}

	// viper.Set("exporter", true)

	if viper.GetBool("exporter") {
		conf.Listener = func(e *types.Event, adapter ports.Adapter) {
			event.Export(e)
		}
	} else {
		conf.Listener = func(e *types.Event, adapter ports.Adapter) {
			handle(e)
		}
	}
	ports.RegisterAdapter(conf)
	log.Infof("Subscribe successful - remote: %v, subscribe: %v", remote, query)

	cancel := func() {
	}
	return cancel, nil
}

func handle(e *types.Event) {
	if e.Source == nil {
		log.Errorf("event's source is nil: %s", e.NodeAddress)
		return
	}
	tags := e.Source.Events
	if tags == nil || len(tags) == 0 {
		log.Errorf("empty event's tags: %s", e.NodeAddress)
		return
	}

	log.Debug("Got event")
}
