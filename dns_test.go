package main

import (
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

type Asd struct {
}

func TestDNSResponse(t *testing.T) {
	lookup := &ContainerLookup{Finder: fakeFinder()}
	r := &dns.Msg{}
	r.SetQuestion("123123123123123.container", dns.TypeTXT)

	d := &DNS{Domain: dns.Fqdn("container")}
	d.ContainerLookup(lookup.FindContainer)

	m := d.BuildResponse(r)

	assert.Equal(t, 1, len(m.Answer))
}
