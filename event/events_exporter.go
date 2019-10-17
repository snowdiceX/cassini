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
	switch tags["message.action"][0] {
	case "transfer":
		checkQos("receive", tags["receive.qos"])
		checkQos("send", tags["send.qos"])
		checkQsc("receive", tags["receive.qscs"])
		checkQsc("send", tags["send.qscs"])
	}
}

func checkQsc(transfer string, vals []string) {
	qscs := checkQscMax(vals)
	for _, c := range qscs {
		if !c.GetAmount().IsInt64() {
			log.Errorf("event qsc amount is not int64: %v", c)
			continue
		}
		if c.GetAmount().Int64() > 0 {
			exporter.Set(exporter.KeyTxMax,
				float64(c.GetAmount().Int64()),
				transfer, c.GetName())
		}
	}
	return
}

func checkQscMax(vals []string) (ret qostypes.QSCs) {
	var err error
	var qscs qostypes.QSCs
	for _, val := range vals {
		if _, qscs, err = qostypes.ParseCoins(val); err != nil {
			log.Errorf("event parse error: %v", err)
		}
		if len(qscs) > 0 {
			ret = append(ret, qscs...)
		}
	}
	return
}

func checkQos(transfer string, vals []string) {
	max := checkQosMax(vals)
	exporter.Set(exporter.KeyTxMax, float64(max), transfer, "qos")
	return
}

func checkQosMax(vals []string) (max int64) {
	var err error
	var tmp int64
	for _, val := range vals {
		if tmp, err = strconv.ParseInt(val, 10, 64); err != nil {
			log.Errorf("event parse error: %v", err)
			return
		}
		if max < tmp {
			max = tmp
		}
	}
	return
}
