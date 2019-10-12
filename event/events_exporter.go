package event

import (
	"strconv"

	"github.com/QOSGroup/cassini/log"
	exporter "github.com/QOSGroup/cassini/prometheus"
	"github.com/QOSGroup/cassini/types"
)

// Import an event for prometheus exporter
func Import(event *types.Event) {
	if event.Source == nil {
		log.Errorf("event's source is nil: %s", event.NodeAddress)
		return
	}
	tags := event.Source.Events
	if tags == nil || len(tags) == 0 {
		log.Errorf("empty event's tags: %s", event.NodeAddress)
		return
	}
	v, err := checkMax(tags["receive.qos"])
	if err != nil {
		log.Warnf("event parse error: %v", err)
		return
	}
	log.Debugf("event receive.qos: %d", v)
	exporter.Set(exporter.KeyTxMax, float64(v), "receive", "qos")
	v, err = checkMax(tags["send.qos"])
	if err != nil {
		log.Warnf("event parse error: %v", err)
		return
	}
	log.Debugf("event send.qos: %d", v)
	exporter.Set(exporter.KeyTxMax, float64(v), "send", "qos")
}

func checkMax(vals []string) (max int64, err error) {
	var tmp int64
	for _, val := range vals {
		if tmp, err = strconv.ParseInt(val, 10, 64); err != nil {
			return
		}
		if max < tmp {
			max = tmp
		}
	}
	return
}
