package main

import (
	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"
	"github.com/prometheus/client_golang/prometheus"
)

type InterfaceCollector struct {
	node *goeapi.Node

	// Gauges
	up         *prometheus.Desc
	info       *prometheus.Desc
	mtuBytes   *prometheus.Desc
	speedBytes *prometheus.Desc

	// Counters
	receiveBytes      *prometheus.Desc
	transmitBytes     *prometheus.Desc
	receivePackets    *prometheus.Desc
	transmitPackets   *prometheus.Desc
	receiveErrs       *prometheus.Desc
	transmitErrs      *prometheus.Desc
	receiveDrop       *prometheus.Desc
	transmitDrop      *prometheus.Desc
	receiveMulticast  *prometheus.Desc
	receiveBroadcast  *prometheus.Desc
	transmitMulticast *prometheus.Desc
	transmitBroadcast *prometheus.Desc
	carrierChanges    *prometheus.Desc
}

func NewInterfaceCollector(node *goeapi.Node) *InterfaceCollector {
	return &InterfaceCollector{
		node: node,

		up: prometheus.NewDesc(
			"node_network_up",
			"Whether the interface is connected.",
			[]string{"device"}, nil,
		),
		info: prometheus.NewDesc(
			"node_network_info",
			"Non-numeric interface metadata.",
			[]string{"device", "operstate", "address", "description", "hardware"}, nil,
		),
		mtuBytes: prometheus.NewDesc(
			"node_network_mtu_bytes",
			"MTU of the interface.",
			[]string{"device"}, nil,
		),
		speedBytes: prometheus.NewDesc(
			"node_network_speed_bytes",
			"Speed of the interface in bytes per second.",
			[]string{"device"}, nil,
		),

		receiveBytes: prometheus.NewDesc(
			"node_network_receive_bytes_total",
			"Total bytes received.",
			[]string{"device"}, nil,
		),
		transmitBytes: prometheus.NewDesc(
			"node_network_transmit_bytes_total",
			"Total bytes transmitted.",
			[]string{"device"}, nil,
		),
		receivePackets: prometheus.NewDesc(
			"node_network_receive_packets_total",
			"Total packets received.",
			[]string{"device"}, nil,
		),
		transmitPackets: prometheus.NewDesc(
			"node_network_transmit_packets_total",
			"Total packets transmitted.",
			[]string{"device"}, nil,
		),
		receiveErrs: prometheus.NewDesc(
			"node_network_receive_errs_total",
			"Total receive errors.",
			[]string{"device"}, nil,
		),
		transmitErrs: prometheus.NewDesc(
			"node_network_transmit_errs_total",
			"Total transmit errors.",
			[]string{"device"}, nil,
		),
		receiveDrop: prometheus.NewDesc(
			"node_network_receive_drop_total",
			"Total received packets dropped.",
			[]string{"device"}, nil,
		),
		transmitDrop: prometheus.NewDesc(
			"node_network_transmit_drop_total",
			"Total transmitted packets dropped.",
			[]string{"device"}, nil,
		),
		receiveMulticast: prometheus.NewDesc(
			"node_network_receive_multicast_total",
			"Total multicast packets received.",
			[]string{"device"}, nil,
		),
		receiveBroadcast: prometheus.NewDesc(
			"node_network_receive_broadcast_total",
			"Total broadcast packets received.",
			[]string{"device"}, nil,
		),
		transmitMulticast: prometheus.NewDesc(
			"node_network_transmit_multicast_total",
			"Total multicast packets transmitted.",
			[]string{"device"}, nil,
		),
		transmitBroadcast: prometheus.NewDesc(
			"node_network_transmit_broadcast_total",
			"Total broadcast packets transmitted.",
			[]string{"device"}, nil,
		),
		carrierChanges: prometheus.NewDesc(
			"node_network_carrier_changes_total",
			"Total carrier link status changes.",
			[]string{"device"}, nil,
		),
	}
}

func (c *InterfaceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.info
	ch <- c.mtuBytes
	ch <- c.speedBytes
	ch <- c.receiveBytes
	ch <- c.transmitBytes
	ch <- c.receivePackets
	ch <- c.transmitPackets
	ch <- c.receiveErrs
	ch <- c.transmitErrs
	ch <- c.receiveDrop
	ch <- c.transmitDrop
	ch <- c.receiveMulticast
	ch <- c.receiveBroadcast
	ch <- c.transmitMulticast
	ch <- c.transmitBroadcast
	ch <- c.carrierChanges
}

func (c *InterfaceCollector) Collect(ch chan<- prometheus.Metric) {
	show := module.Show(c.node)
	result := show.ShowInterfaces()

	for name, iface := range result.Interfaces {
		up := 0.0
		if iface.InterfaceStatus == "connected" {
			up = 1.0
		}
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, up, name)

		ch <- prometheus.MustNewConstMetric(c.info, prometheus.GaugeValue, 1,
			name, iface.InterfaceStatus, iface.PhysicalAddress, iface.Description, iface.Hardware,
		)

		ch <- prometheus.MustNewConstMetric(c.mtuBytes, prometheus.GaugeValue, float64(iface.Mtu), name)
		ch <- prometheus.MustNewConstMetric(c.speedBytes, prometheus.GaugeValue, float64(iface.Bandwidth)/8, name)

		counters := iface.InterfaceCounters
		ch <- prometheus.MustNewConstMetric(c.receiveBytes, prometheus.CounterValue, float64(counters.InOctets), name)
		ch <- prometheus.MustNewConstMetric(c.transmitBytes, prometheus.CounterValue, float64(counters.OutOctets), name)

		inPkts := counters.InUcastPkts + counters.InMulticastPkts + counters.InBroadcastPkts
		outPkts := counters.OutUcastPkts + counters.OutMulticastPkts + counters.OutBroadcastPkts
		ch <- prometheus.MustNewConstMetric(c.receivePackets, prometheus.CounterValue, float64(inPkts), name)
		ch <- prometheus.MustNewConstMetric(c.transmitPackets, prometheus.CounterValue, float64(outPkts), name)

		ch <- prometheus.MustNewConstMetric(c.receiveErrs, prometheus.CounterValue, float64(counters.TotalInErrors), name)
		ch <- prometheus.MustNewConstMetric(c.transmitErrs, prometheus.CounterValue, float64(counters.TotalOutErrors), name)
		ch <- prometheus.MustNewConstMetric(c.receiveDrop, prometheus.CounterValue, float64(counters.InDiscards), name)
		ch <- prometheus.MustNewConstMetric(c.transmitDrop, prometheus.CounterValue, float64(counters.OutDiscards), name)
		ch <- prometheus.MustNewConstMetric(c.receiveMulticast, prometheus.CounterValue, float64(counters.InMulticastPkts), name)
		ch <- prometheus.MustNewConstMetric(c.receiveBroadcast, prometheus.CounterValue, float64(counters.InBroadcastPkts), name)
		ch <- prometheus.MustNewConstMetric(c.transmitMulticast, prometheus.CounterValue, float64(counters.OutMulticastPkts), name)
		ch <- prometheus.MustNewConstMetric(c.transmitBroadcast, prometheus.CounterValue, float64(counters.OutBroadcastPkts), name)
		ch <- prometheus.MustNewConstMetric(c.carrierChanges, prometheus.CounterValue, float64(counters.LinkStatusChanges), name)
	}
}
