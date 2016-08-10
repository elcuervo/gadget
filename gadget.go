package main

import (
	"flag"
	"log"
	"os"
)

var addr = flag.String("address", ":53", "Address to bind.")
var domain = flag.String("domain", "container", "Domain to use for lookups.")
var socket = flag.String("socket", "/var/run/docker.sock", "Docker socket")

func main() {
	flag.Parse()

	if _, err := os.Stat(*socket); os.IsNotExist(err) {
		log.Fatal("I need access to the Docker socket to be able to work.")
	}

	server := NewDNSServer(*addr, *domain)
	lookup := NewContainerLookup(*socket)

	server.ContainerLookup(lookup.FindContainer)

	log.Printf("Serving %s and with FQDN: %s", *addr, *domain)

	go server.Serve()
	defer server.Shutdown()

	server.Wait()
}
