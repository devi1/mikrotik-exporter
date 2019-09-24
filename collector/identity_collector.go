package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/routeros.v2/proto"
)

type identityCollector struct {
	identityDesc *prometheus.Desc
}

func (c *identityCollector) init() {
	const prefix = "system"

	labelNames := []string{"name", "address", "identity"}
	c.identityDesc = description(prefix, "identity", "RouterOS system identity", labelNames)
}

func newIdentityCollector() routerOSCollector {
	c := &identityCollector{}
	c.init()
	return c
}

func (c *identityCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.identityDesc
}

func (c *identityCollector) collect(ctx *collectorContext) error {
	stats, err := c.fetch(ctx)
	if err != nil {
		return err
	}

	for _, re := range stats {
		c.collectMetric(ctx, re)
	}

	return nil
}

func (c *identityCollector) fetch(ctx *collectorContext) ([]*proto.Sentence, error) {
	reply, err := ctx.client.Run("/system/identity/print", "=.proplist=name")
	if err != nil {
		log.WithFields(log.Fields{
			"device": ctx.device.Name,
			"error":  err,
		}).Error("error fetching identity name")
		return nil, err
	}

	return reply.Re, nil
}

func (c *identityCollector) collectMetric(ctx *collectorContext, re *proto.Sentence) {
	v := 1.0
	identity := re.Map["name"]

	ctx.ch <- prometheus.MustNewConstMetric(c.identityDesc, prometheus.CounterValue, v, ctx.device.Name, ctx.device.Address, identity)
}
