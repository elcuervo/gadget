package main

import (
	"flag"
	"log"
)

var addr = flag.String("address", ":53", "Address to bind.")
var domain = flag.String("domain", "container", "Domain to use for lookups.")
var socket = flag.String("socket", "/var/run/docker.sock", "Docker socket")

func main() {
	flag.Parse()

	server := NewDnsServer(*addr, *domain)

	log.Printf("Serving %s and with FQDN: %s", *addr, *domain)

	go server.Serve()
	defer server.Shutdown()
	server.Wait()
}
