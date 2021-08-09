package collectors

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rocket-pool/rocketpool-go/deposit"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

const namespace = "rocketpool"


// Represents the collector for the Demand metrics
type DemandCollector struct {
	// The amount of ETH currently in the Deposit Pool
	depositPoolBalance 			*prometheus.Desc

	// The excess ETH balance of the Deposit Pool
	depositPoolExcess			*prometheus.Desc

	// The total ETH capacity of the Minipool queue
	totalMinipoolCapacity 		*prometheus.Desc

	// The effective ETH capacity of the Minipool queue
	effectiveMinipoolCapacity 	*prometheus.Desc

	// The Rocket Pool contract manager
	rp 							*rocketpool.RocketPool
}


// Create a new DemandCollector instance
func NewDemandCollector(rp *rocketpool.RocketPool) *DemandCollector {
	subsystem := "demand"
	return &DemandCollector{
		depositPoolBalance: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "deposit_pool_balance"),
			"The amount of ETH currently in the Deposit Pool",
			nil, nil,
		),
		depositPoolExcess: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "deposit_pool_excess"),
			"The excess ETH balance of the Deposit Pool",
			nil, nil,
		),
		totalMinipoolCapacity: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "total_minipool_capacity"),
			"The total ETH capacity of the Minipool queue",
			nil, nil,
		),
		effectiveMinipoolCapacity: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "effective_minipool_capacity"),
			"The effective ETH capacity of the Minipool queue",
			nil, nil,
		),
		rp: rp,
	}
}


// Write metric descriptions to the Prometheus channel
func (collector *DemandCollector) Describe(channel chan<- *prometheus.Desc) {
	channel <- collector.depositPoolBalance
	channel <- collector.depositPoolExcess
	channel <- collector.totalMinipoolCapacity
	channel <- collector.effectiveMinipoolCapacity
}


// Collect the latest metric values and pass them to Prometheus
func (collector *DemandCollector) Collect(channel chan<- prometheus.Metric) {
 
	// Get the Deposit Pool balance
	depositPoolBalance, err := deposit.GetBalance(collector.rp, nil)
	balanceFloat := float64(-1)
	if err != nil {
		log.Printf("Error getting deposit pool balance: %s", err)
	} else {
		balanceFloat = eth.WeiToEth(depositPoolBalance)
	}
    channel <- prometheus.MustNewConstMetric(
		collector.depositPoolBalance, prometheus.GaugeValue, balanceFloat)
	
	// Get the deposit pool excess
	depositPoolExcess, err := deposit.GetExcessBalance(collector.rp, nil)
	excessFloat := float64(-1)
	if err != nil {
		log.Printf("Error getting deposit pool excess: %s", err)
	} else {
		excessFloat = eth.WeiToEth(depositPoolExcess)
	}
    channel <- prometheus.MustNewConstMetric(
		collector.depositPoolExcess, prometheus.GaugeValue, excessFloat)


	// Get the total and effective minipool capacities
	minipoolQueueCapacity, err := minipool.GetQueueCapacity(collector.rp, nil)
	totalFloat := float64(-1) 
	effectiveFloat := float64(-1)
	if err != nil {
		log.Printf("Error getting minipool queue capacity: %s", err)
	} else {
		totalFloat = eth.WeiToEth(minipoolQueueCapacity.Total)
		effectiveFloat = eth.WeiToEth(minipoolQueueCapacity.Effective)
	}
    channel <- prometheus.MustNewConstMetric(
		collector.totalMinipoolCapacity, prometheus.GaugeValue, totalFloat)
	channel <- prometheus.MustNewConstMetric(
		collector.effectiveMinipoolCapacity, prometheus.GaugeValue, effectiveFloat)
}