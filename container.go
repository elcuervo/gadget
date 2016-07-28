package main

import (
	"errors"
	"net"
	"strconv"

	"github.com/docker/engine-api/client"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

type ContainerInfo struct {
	ID     string
	Image  string
	Name   string
	Status string
	Fqdn   string
}

type ContainerService struct {
	Addr string
	Port int
}

type Container struct {
	Info     *ContainerInfo
	IPs      []net.IP
	Services []ContainerService
}

type ContainerLookup struct {
	docker *client.Client
}

func (c *Container) ToTXT() []string {
	return []string{
		"id=" + c.Info.ID,
		"name=" + c.Info.Name,
		"image=" + c.Info.Image,
		"status=" + c.Info.Status,
	}
}

func (l *ContainerLookup) Find(id string) (Container, error) {
	var ips []net.IP
	var services []ContainerService

	i, err := l.docker.ContainerInspect(context.Background(), id)

	if err != nil {
		return Container{}, errors.New("Container not found")
	}

	info := &ContainerInfo{
		ID:     i.ID[:12],
		Image:  i.Config.Image,
		Name:   i.Name[1:],
		Status: i.State.Status,
		Fqdn:   id + "." + dns.Fqdn(*domain),
	}

	for _, netw := range i.NetworkSettings.Networks {
		ips = append(ips, net.ParseIP(netw.IPAddress))
	}

	for cport, hports := range i.NetworkSettings.Ports {
		port := cport.Int()

		services = append(services, ContainerService{info.Fqdn, port})

		if len(hports) == 0 {
			services = append(services, ContainerService{"localhost.localdomain.", 0})
		} else {
			for _, hport := range hports {
				hostIP := net.ParseIP(hport.HostIP)

				if hostIP.Equal(net.IPv4(0, 0, 0, 0)) {
					hostIP = net.IPv4(127, 0, 0, 1)
				}

				hosts, _ := net.LookupAddr(hostIP.String())
				port, _ := strconv.Atoi(hport.HostPort)

				services = append(services, ContainerService{hosts[0], port})
			}
		}
	}

	return Container{
		Info:     info,
		IPs:      ips,
		Services: services,
	}, nil
}

func NewContainerLookup(addr string) *ContainerLookup {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(addr, "v1.22", nil, defaultHeaders)

	if err != nil {
		panic(err)
	}

	l := new(ContainerLookup)
	l.docker = cli

	return l
}

func FindContainer(containerID string) (Container, error) {
	lookup := NewContainerLookup("unix://" + *socket)
	container, err := lookup.Find(containerID)

	if err != nil {
		return Container{}, errors.New("Container " + containerID + " not found")
	} else {
		return container, nil
	}
}
