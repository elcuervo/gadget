package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/engine-api/client"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

var addr = flag.String("address", ":53", "Address to bind.")
var domain = flag.String("domain", "container", "Domain to use for lookups.")

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

func containerLookup(containerID string) (Container, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
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
			Fqdn:  containerID + ".container.",
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

func buildResponse(r *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	for _, q := range m.Question {
		var rrs []dns.RR

		log.Printf("Answering %s type: %d", q.Name, q.Qtype)

		host := q.Name[0 : len(q.Name)-1]
		containerID := strings.TrimSuffix(host, filepath.Ext(host))
		container, _ := containerLookup(containerID)

		log.Print("Building TXT record")
		txt := new(dns.TXT)
		txt.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 3600}
		txt.Txt = container.ToTXT()

		rrs = append(rrs, txt)

		log.Print("Building A record")
		for _, ip := range container.IPs {
			a := new(dns.A)
			a.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 3600}
			a.A = ip

			rrs = append(rrs, a)
		}

		log.Print("Building SRV record")
		for i, service := range container.Services {
			srv := new(dns.SRV)
			srv.Hdr = dns.RR_Header{q.Name, dns.TypeSRV, dns.ClassINET, 3600, 0}
			srv.Port = uint16(service.Port)
			srv.Weight = uint16(i)
			srv.Target = service.Addr

			rrs = append(rrs, srv)
		}

		m.Answer = append(m.Answer, rrs...)
	}

	return m
}

func dnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	log.Printf("Incoming DNS request")

	m := buildResponse(r)

	if r.IsTsig() != nil {
		if w.TsigStatus() == nil {
			m.SetTsig("axfr.", dns.HmacMD5, 300, time.Now().Unix())
		}
	}

	w.WriteMsg(m)

}

func main() {
	flag.Parse()
	addr := *addr
	domain := *domain
	prefix := dns.Fqdn(domain)

	server := &dns.Server{Addr: addr, Net: "udp"}
	go server.ListenAndServe()

	dns.HandleFunc(prefix, dnsRequest)

	defer server.Shutdown()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Serving %s and with FQDN: %s", addr, prefix)

	for {
		select {
		case <-sig:
			log.Printf("Shutting Down")
			os.Exit(0)
		}
	}
}
