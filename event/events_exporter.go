package event

import (
	"strconv"

	"github.com/QOSGroup/cassini/log"
	exporter "github.com/QOSGroup/cassini/prometheus"
	"github.com/QOSGroup/cassini/types"
	qostypes "github.com/QOSGroup/qos/types"
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
	checkQosMax("receive", tags["receive.qos"])
	checkQosMax("send", tags["send.qos"])
	checkQscMax("receive", tags["receive.qscs"])
	checkQscMax("send", tags["send.qscs"])
}

func checkQscMax(transfer string, vals []string) {
	var err error
	// var tmp btypes.BigInt
	var qscs qostypes.QSCs
	for _, val := range vals {
		if _, qscs, err = qostypes.ParseCoins(val); err != nil {
			log.Errorf("event parse error: %v", err)
		}
		for _, c := range qscs {
			if c.GetAmount().IsInt64() {
				exporter.Set(exporter.KeyTxMax,
					float64(c.GetAmount().Int64()),
					transfer, c.GetName())
			} else {
				log.Errorf("event qsc amount is not int64: %v", c)
			}
		}
	}
	return
}

func checkQosMax(transfer string, vals []string) {
	var err error
	var max, tmp int64
	for _, val := range vals {
		if tmp, err = strconv.ParseInt(val, 10, 64); err != nil {
			log.Errorf("event parse error: %v", err)
			return
		}
		if max < tmp {
			max = tmp
		}
	}
	exporter.Set(exporter.KeyTxMax, float64(max), transfer, "qos")
	return
}
