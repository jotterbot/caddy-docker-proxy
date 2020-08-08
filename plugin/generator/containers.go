package generator

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/lucaslorentz/caddy-docker-proxy/plugin/v2/caddyfile"
)

func (g *CaddyfileGenerator) getContainerCaddyfile(container *types.Container, logsBuffer *bytes.Buffer) (*caddyfile.Container, error) {
	caddyLabels := g.filterLabels(container.Labels)

	return labelsToCaddyfile(caddyLabels, container, func() ([]string, error) {
		return g.getContainerPublicPort(container, logsBuffer)
	})
}

func (g *CaddyfileGenerator) getContainerIPAddresses(container *types.Container, logsBuffer *bytes.Buffer, ingress bool) ([]string, error) {
	ips := []string{}

	for _, network := range container.NetworkSettings.Networks {
		if !ingress || g.ingressNetworks[network.NetworkID] {
			ips = append(ips, network.IPAddress)
		}
	}

	if len(ips) == 0 {
		logsBuffer.WriteString(fmt.Sprintf("[WARNING] Container %v and caddy are not in same network\n", container.ID))
	}

	return ips, nil
}

func (g *CaddyfileGenerator) getContainerPublicPort(container *types.Container, logsBuffer *bytes.Buffer) ([]string, error) {
	ports := []string{}
	for _, port := range container.Ports {
			public_port := strconv.Itoa(int(port.PublicPort))
			ports = append(ports, public_port)
	}
	return ports,nil
}