package main

import (
	"log"
	"net"
	"strconv"

	"github.com/docker/engine-api/client"
	"golang.org/x/net/context"
)

type ContainerInfo struct {
	ID    string
	Image string
	Name  string
	Fqdn  string
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

func (c *Container) ToTXT() []string {
	var txt []string

	txt = append(txt, "id="+c.Info.ID)
	txt = append(txt, "name="+c.Info.Name)
	txt = append(txt, "image="+c.Info.Image)

	return txt
}

func FindContainer(containerID string) (Container, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix://"+*socket, "v1.22", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}
	containerInspect, err := cli.ContainerInspect(context.Background(), containerID)

	if err != nil {
		panic(err)
	} else {
		var ips []net.IP
		var services []ContainerService

		//log.Printf("%+v", containerInspect.Config)

		info := &ContainerInfo{
			ID:    containerInspect.ID[:12],
			Image: containerInspect.Config.Image,
			Name:  containerInspect.Name[1:],
			Fqdn:  containerID + *domain,
		}

		for _, netw := range containerInspect.NetworkSettings.Networks {
			ips = append(ips, net.ParseIP(netw.IPAddress))
		}

		for cport, hports := range containerInspect.NetworkSettings.Ports {
			port := cport.Int()

			services = append(services, ContainerService{info.Fqdn, port})

			log.Printf("%+v", hports)

			if len(hports) == 0 {
				services = append(services, ContainerService{"localhost.localdomain.", 0})
			}

			for _, hport := range hports {
				hostIP := net.ParseIP(hport.HostIP)

				if hostIP.Equal(net.IPv4(0, 0, 0, 0)) {
					hostIP = net.IPv4(127, 0, 0, 1)
				} else {
					panic("god")
				}

				hosts, _ := net.LookupAddr(hostIP.String())
				port, _ := strconv.Atoi(hport.HostPort)
				services = append(services, ContainerService{hosts[0], port})
			}
		}

		return Container{
			Info:     info,
			IPs:      ips,
			Services: services,
		}, nil
	}
}
