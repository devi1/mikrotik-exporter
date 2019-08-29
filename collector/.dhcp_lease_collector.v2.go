package collector

import (
//	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/routeros.v2/proto"
)

type dhcpLeaseCollector struct {
	props        []string
	descriptions *prometheus.Desc
}

func (c *dhcpLeaseCollector) init() {
	c.props = []string{"active-mac-address", "status", "expires-after", "active-address", "host-name"}
	labelNames := []string{"name", "address", "activemacaddress", "status", "expiresafter", "activeaddress", "hostname"}
	c.descriptions = description("dhcp", "leases_metrics", "DHCP Leases mertics", labelNames)
}

func newDHCPLCollector() routerOSCollector {
	c := &dhcpLeaseCollector{}
	c.init()
	return c
}

func (c *dhcpLeaseCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.descriptions
}

func (c *dhcpLeaseCollector) collect(ctx *collectorContext) error {
	names, err := c.fetch(ctx)
	if err != nil {
		return err
	}

	for _, n := range names {
		c.collectForDHCPLease(re, ctx)
	}

	return nil
}

func (c *dhcpLeaseCollector) fetch(ctx *collectorContext) ([]*proto.Sentence, error) {
	reply, err := ctx.client.Run("/ip/dhcp-server/lease/print", "=.proplist="+strings.Join(c.props, ","))
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching DHCP leases metrics")
		return nil, err
	}

	return reply.Re, nil
}



func (c *dhcpLeaseCollector) collectForDHCPLease(re *proto.Sentence, ctx *collectorContext) {
	reply, err := ctx.client.Run("/ip/dhcp-server/lease/print", "=.proplist="+strings.Join(c.props, ","))
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching DHCP lease metrics")
		return err
	}

	v, err := strconv.ParseFloat(reply.Done.Map["ret"], 32)
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error parsing DHCP lease metrics")
		return err
	}

	ctx.ch <- prometheus.MustNewConstMetric(c.descriptions, prometheus.GaugeValue, v, ctx.device.Name, ctx.device.Address, activemacaddress, status, expiresafter, activeaddress, hostname)
	return nil
}

//func (c *dhcpLeaseCollector) fetchDHCPServerNames(ctx *collectorContext) ([]string, error) {
//	reply, err := ctx.client.Run("/ip/dhcp-server/print", "=.proplist=name")
//	if err != nil {
//		log.WithFields(log.Fields{
//			"device": ctx.device.Name,
//			"error":  err,
//		}).Error("error fetching DHCP server names")
//		return nil, err
//	}
//
//	names := []string{}
//	for _, re := range reply.Re {
//		names = append(names, re.Map["name"])
//	}
//
//	return names, nil
//}

/*
func (c *dhcpLeaseCollector) collectForStat(re *proto.Sentence, ctx *collectorContext) {
	activemacaddress := re.Map["active-mac-address"]
	status := re.Map["status"]
	expiresafter := re.Map["expires-after"]
	activeaddress := re.Map["active-address"]
	hostname := re.Map["host-name"]

	//for _, p := range c.props {
	c.collectMetricForProperty(activemacaddress, status, expiresafter, activeaddress, hostname, re, ctx)
	//}
}
*/

/*
func (c *dhcpLeaseCollector) collectMetricForProperty(activemacaddress, status, expiresafter, activeaddress, hostname string, re *proto.Sentence, ctx *collectorContext) {
	//desc := c.descriptions[property]
	//v, err := c.parseValueForProperty(property, re.Map[property])
	//desc := "somedescr"
	labelNames := []string{"name", "address", "server"}
	desc := description("prefix", "leases_active_count", "number of active leases per DHCP server", labelNames)
	v := 1.0
	err := 0
	//	if value := re.Map[property]; value != "" {
	//		v, err := strconv.ParseFloat(value, 64)
	if err != 0 {
		log.WithFields(log.Fields{
			"device":           ctx.device.Name,
			"activemacaddress": activemacaddress,
			"status":           status,
			"expiresafter":     expiresafter,
			"activeaddress":    activeaddress,
			"hostname":         hostname,
			//"property":         property,
			//"value": v,
			"error": err,
		}).Error("error parsing DHCP Leases metric value")
		return
	}
	ctx.ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, v, ctx.device.Name, ctx.device.Address, activemacaddress, status, expiresafter, activeaddress, hostname)
	//}
}
*/
