package main

import (
	"context"
	"fmt"

	"github.com/huangdao/cassini/config"
	"github.com/huangdao/cassini/event"
	"github.com/huangdao/cassini/log"

	tmtypes "github.com/tendermint/tendermint/types"
)

// 命令行 events 命令执行方法
var events = func(conf *config.Config) (context.CancelFunc, error) {
	cancelFunc, err := Subscribe(conf.EventsListen, conf.EventsQuery)
	if err != nil {
		return nil, err
	}
	cancel := func() {
		cancelFunc()
		log.Debug("Cancel events subscribe service")
	}
	return cancel, nil
}

//Subscribe 从websocket服务端订阅event
//remote 服务端地址 example  "tcp://127.0.0.1:27657"
func Subscribe(remote string, query string) (context.CancelFunc, error) {
	txs := make(chan interface{})
	cancel, err := event.SubscribeRemote(remote, "cassini-events", query, txs)
	if err != nil {
		log.Errorf("Remote [%s] : '%s'\n", remote, err)
		return nil, err
	}
	go func() {
		for e := range txs {
			fmt.Println("Got Tx event - ", e.(tmtypes.EventDataTx)) //注：e类型断言为types.CassiniEventDataTx 类型
			for _, tto := range e.(tmtypes.EventDataTx).Result.Tags {
				kv := tto //interface{}(tto).(common.KVPair)
				fmt.Println(string(kv.Key), string(kv.Value))
			}
		}
	}()
	return cancel, nil
}
