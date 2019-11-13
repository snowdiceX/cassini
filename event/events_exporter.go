package event

import (
	"strconv"

	"github.com/QOSGroup/cassini/log"
	"github.com/QOSGroup/cassini/types"
	qostypes "github.com/QOSGroup/qos/types"
	"github.com/snowdiceX/exporter"
	"github.com/spf13/viper"
)

// nolint
const (
	MetricAddr      = ":39099"
	MetricPath      = "/metrics"
	KeyDuration     = "duration"
	KeyPrefix       = "cassini_"
	KeyQueueSize    = "queue_size"
	KeyQueue        = "queue"
	KeyAdaptors     = "adaptors"
	KeyTxMax        = "tx_max"
	KeyTxsWait      = "txs_wait"
	KeyTxCost       = "tx_cost"
	KeyTxsPerSecond = "txs_per_second"
	KeyErrors       = "errors"
)

func init() {
	initMcs := []*exporter.MetricConfig{
		&exporter.MetricConfig{
			Key:    KeyQueueSize,
			Type:   "ImmutableGaugeMetric",
			Help:   "Size of queue",
			Labels: []string{"type"}},
		&exporter.MetricConfig{
			Key:  KeyErrors,
			Type: "CounterMetric",
			Help: "Count of running errors"},
		&exporter.MetricConfig{
			Key:    KeyTxMax,
			Type:   "TxMaxGaugeMetric",
			Help:   "Max value of transfer txs per minute",
			Labels: []string{"transfer", "token", "txhash"}},
		&exporter.MetricConfig{
			Key:  KeyTxsPerSecond,
			Type: "TickerGaugeMetric",
			Help: "Number of relayed tx per second"},
		&exporter.MetricConfig{
			Key:  KeyTxCost,
			Type: "TickerGaugeMetric",
			Help: "Time(milliseconds) cost of lastest tx relay"}}

	viper.Set(exporter.KeyMetricAddr, MetricAddr)
	viper.Set(exporter.KeyMetricPath, MetricPath)
	viper.Set(exporter.KeyMetricPrefix, KeyPrefix)

	viper.Set(exporter.KeyMetricType, initMcs)

	viper.Set(KeyDuration, "30")
}

// Export an event for prometheus exporter
func Export(event *types.Event) {
	if event.Source == nil {
		log.Errorf("event's source is nil: %s", event.NodeAddress)
		return
	}
	tags := event.Source.Events
	if tags == nil || len(tags) == 0 {
		log.Errorf("empty event's tags: %s", event.NodeAddress)
		return
	}
	action := tags["message.action"]
	if action == nil || len(action) == 0 {
		log.Errorf("no message.action tag: %s", event.NodeAddress)
		return
	}
	switch action[0] {
	case "transfer":
		txhash := tags["tx.hash"]
		hash := ""
		if len(txhash) > 0 {
			hash = txhash[0]
		}
		checkQos("receive", tags["receive.qos"], hash)
		checkQos("send", tags["send.qos"], hash)
		checkQsc("receive", tags["receive.qscs"], hash)
		checkQsc("send", tags["send.qscs"], hash)
	}
}

func checkQsc(transfer string, vals []string, hash string) {
	qscs := checkQscMax(vals)
	for _, c := range qscs {
		if !c.GetAmount().IsInt64() {
			log.Errorf("event qsc amount is not int64: %v", c)
			continue
		}
		if c.GetAmount().Int64() > 0 {
			exporter.Set(KeyTxMax,
				float64(c.GetAmount().Int64()),
				transfer, c.GetName(), hash)
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

func checkQos(transfer string, vals []string, hash string) {
	max := checkQosMax(vals)
	exporter.Set(KeyTxMax, float64(max), transfer, "qos", hash)
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
