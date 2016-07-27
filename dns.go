package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/miekg/dns"
)

type Dns struct {
	Domain string
	server *dns.Server
}

type DnsQueryResponder struct {
	Question  dns.Question
	Container Container
}

func (r *DnsQueryResponder) buildHdr(rtype uint16) dns.RR_Header {
	return dns.RR_Header{
		Name:   r.Question.Name,
		Rrtype: rtype,
		Class:  dns.ClassINET,
		Ttl:    3600,
	}
}

func (r *DnsQueryResponder) buildRR(rtype string) []dns.RR {
	var rrs []dns.RR

	switch rtype {
	case "TXT":
		rr := new(dns.TXT)
		rr.Hdr = r.buildHdr(dns.TypeTXT)
		rr.Txt = r.Container.ToTXT()

		rrs = append(rrs, rr)
	case "A":
		for _, ip := range r.Container.IPs {
			rr := new(dns.A)
			rr.Hdr = r.buildHdr(dns.TypeA)
			rr.A = ip

			rrs = append(rrs, rr)
		}

	case "SRV":
		for i, service := range r.Container.Services {
			rr := new(dns.SRV)
			rr.Hdr = r.buildHdr(dns.TypeSRV)
			rr.Port = uint16(service.Port)
			rr.Weight = uint16(i)
			rr.Target = service.Addr

			rrs = append(rrs, rr)
		}
	}

	return rrs
}

func (d *Dns) buildResponse(r *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	for _, q := range m.Question {
		var rrs []dns.RR

		log.Printf("Answering %s type: %d", q.Name, q.Qtype)

		host := q.Name[0 : len(q.Name)-1]
		containerID := strings.TrimSuffix(host, filepath.Ext(host))
		container, _ := FindContainer(containerID)
		res := &DnsQueryResponder{q, container}

		rrs = append(rrs, res.buildRR("TXT")...)
		rrs = append(rrs, res.buildRR("A")...)
		rrs = append(rrs, res.buildRR("SRV")...)

		m.Answer = append(m.Answer, rrs...)
	}

	return m
}

func (d *Dns) handleDns(w dns.ResponseWriter, r *dns.Msg) {
	m := d.buildResponse(r)

	if r.IsTsig() != nil {
		if w.TsigStatus() == nil {
			m.SetTsig("axfr.", dns.HmacMD5, 300, time.Now().Unix())
		}
	}

	w.WriteMsg(m)

}

func (d *Dns) Wait() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			log.Printf("Shutting Down")
			os.Exit(0)
		}
	}
}

func (d *Dns) Shutdown() {
	d.server.Shutdown()
}

func (d *Dns) Serve() {
	d.server.ListenAndServe()
}

func NewDnsServer(address, domain string) *Dns {
	d := new(Dns)
	d.server = &dns.Server{Addr: address, Net: "udp"}
	d.Domain = dns.Fqdn(domain)

	dns.HandleFunc(d.Domain, d.handleDns)

	return d
}
