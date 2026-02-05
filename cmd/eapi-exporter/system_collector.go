package main

import (
	"log"
	"strings"

	"github.com/aristanetworks/goeapi"
	"github.com/aristanetworks/goeapi/module"
	"github.com/prometheus/client_golang/prometheus"
)

type SystemCollector struct {
	node *goeapi.Node

	// System
	bootTime *prometheus.Desc

	// Power supply
	powerInfo          *prometheus.Desc
	powerStatus        *prometheus.Desc
	powerCapacityWatts *prometheus.Desc
	powerInputAmps     *prometheus.Desc
	powerOutputAmps    *prometheus.Desc
	powerOutputWatts   *prometheus.Desc
	powerUptime        *prometheus.Desc
	powerTempCelsius   *prometheus.Desc
	powerFanSpeed      *prometheus.Desc
	powerFanStatus     *prometheus.Desc
}

func NewSystemCollector(node *goeapi.Node) *SystemCollector {
	return &SystemCollector{
		node: node,

		bootTime: prometheus.NewDesc(
			"node_boot_time_seconds",
			"Unix timestamp of system boot time.",
			nil, nil,
		),

		powerInfo: prometheus.NewDesc(
			"node_power_supply_info",
			"Power supply metadata.",
			[]string{"supply", "model"}, nil,
		),
		powerStatus: prometheus.NewDesc(
			"node_power_supply_status",
			"Power supply status, 1 if ok.",
			[]string{"supply"}, nil,
		),
		powerCapacityWatts: prometheus.NewDesc(
			"node_power_supply_capacity_watts",
			"Power supply capacity in watts.",
			[]string{"supply"}, nil,
		),
		powerInputAmps: prometheus.NewDesc(
			"node_power_supply_input_current_amperes",
			"Power supply input current in amperes.",
			[]string{"supply"}, nil,
		),
		powerOutputAmps: prometheus.NewDesc(
			"node_power_supply_output_current_amperes",
			"Power supply output current in amperes.",
			[]string{"supply"}, nil,
		),
		powerOutputWatts: prometheus.NewDesc(
			"node_power_supply_output_watts",
			"Power supply output power in watts.",
			[]string{"supply"}, nil,
		),
		powerUptime: prometheus.NewDesc(
			"node_power_supply_uptime_seconds",
			"Power supply uptime in seconds.",
			[]string{"supply"}, nil,
		),
		powerTempCelsius: prometheus.NewDesc(
			"node_power_supply_temp_celsius",
			"Power supply temperature in celsius.",
			[]string{"supply", "sensor"}, nil,
		),
		powerFanSpeed: prometheus.NewDesc(
			"node_power_supply_fan_speed",
			"Power supply fan speed.",
			[]string{"supply", "fan"}, nil,
		),
		powerFanStatus: prometheus.NewDesc(
			"node_power_supply_fan_status",
			"Power supply fan status, 1 if ok.",
			[]string{"supply", "fan"}, nil,
		),
	}
}

func (c *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.bootTime
	ch <- c.powerInfo
	ch <- c.powerStatus
	ch <- c.powerCapacityWatts
	ch <- c.powerInputAmps
	ch <- c.powerOutputAmps
	ch <- c.powerOutputWatts
	ch <- c.powerUptime
	ch <- c.powerTempCelsius
	ch <- c.powerFanSpeed
	ch <- c.powerFanStatus
}

func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	show := module.Show(c.node)

	// Boot time from show version
	version := show.ShowVersion()
	if version.BootupTimestamp > 0 {
		ch <- prometheus.MustNewConstMetric(c.bootTime, prometheus.GaugeValue, version.BootupTimestamp)
	}

	// Power supplies from show environment power
	power, err := show.ShowEnvironmentPower()
	if err != nil {
		log.Printf("failed to get power data: %v", err)
		return
	}

	for name, psu := range power.PowerSupplies {
		status := 0.0
		if strings.EqualFold(psu.State, "ok") {
			status = 1.0
		}

		ch <- prometheus.MustNewConstMetric(c.powerInfo, prometheus.GaugeValue, 1, name, psu.ModelName)
		ch <- prometheus.MustNewConstMetric(c.powerStatus, prometheus.GaugeValue, status, name)
		ch <- prometheus.MustNewConstMetric(c.powerCapacityWatts, prometheus.GaugeValue, float64(psu.Capacity), name)
		ch <- prometheus.MustNewConstMetric(c.powerInputAmps, prometheus.GaugeValue, psu.InputCurrent, name)
		ch <- prometheus.MustNewConstMetric(c.powerOutputAmps, prometheus.GaugeValue, psu.OutputCurrent, name)
		ch <- prometheus.MustNewConstMetric(c.powerOutputWatts, prometheus.GaugeValue, psu.OutputPower, name)
		ch <- prometheus.MustNewConstMetric(c.powerUptime, prometheus.GaugeValue, psu.Uptime, name)

		for sensorName, sensor := range psu.TempSensors {
			ch <- prometheus.MustNewConstMetric(c.powerTempCelsius, prometheus.GaugeValue, float64(sensor.Temperature), name, sensorName)
		}

		for fanName, fan := range psu.Fans {
			fanStatus := 0.0
			if strings.EqualFold(fan.Status, "ok") {
				fanStatus = 1.0
			}
			ch <- prometheus.MustNewConstMetric(c.powerFanSpeed, prometheus.GaugeValue, float64(fan.Speed), name, fanName)
			ch <- prometheus.MustNewConstMetric(c.powerFanStatus, prometheus.GaugeValue, fanStatus, name, fanName)
		}
	}
}
