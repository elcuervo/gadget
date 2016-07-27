# Inspector Gadget
_Don't worry Cief, I'm always on duty!_

![](http://clipset.20minutos.es/wp-content/uploads/2013/10/gadget.jpg)

Gadget is a custom DNS server with the whole purpose on helping you find
information about a given container running in your host.

It answers by default to any `.container` and will return general container
information in a `TXT` record, the internal IP address in an `A` record and if
the container has exposed ports you'll see them in the `SRV` records.

## Installation

```bash
go get github.com/elcuervo/gadget
```

## Configuration

You have several flags to configure `gadget`.

```bash
gadget -domain ships -address 10.0.0.1:5353 -socket /var/sockets/docker.sock
```

Being the default values `-domain container`, `-address :53` and `-socket /var/run/docker.sock`

## DNS Querying

```bash
dig [container_id].container @[DNS Address]
```

Note that a container can be found also by it name.

Example:

```bash
$ dig 89532c3f0369.container
```

```bash
; <<>> DiG 9.10.4-P2 <<>> 89532c3f0369.container @127.0.0.1
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 12838
;; flags: qr rd; QUERY: 1, ANSWER: 6, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;89532c3f0369.container.		IN	A

;; ANSWER SECTION:
89532c3f0369.container.	3600	IN	TXT	"id=89532c3f0369" "name=drunk_newton" "image=redis" "status=running"
89532c3f0369.container.	3600	IN	A	172.17.0.3
89532c3f0369.container.	3600	IN	SRV	0 0 6237 89532c3f0369.container.
89532c3f0369.container.	3600	IN	SRV	0 1 32768 localhost.localdomain.
89532c3f0369.container.	3600	IN	SRV	0 2 6379 89532c3f0369.container.
89532c3f0369.container.	3600	IN	SRV	0 3 0 localhost.localdomain.

;; Query time: 5 msec
;; SERVER: 127.0.0.1#53(127.0.0.1)
;; WHEN: Wed Jul 27 17:21:36 UYT 2016
;; MSG SIZE  rcvd: 427
```

### TXT records

They include the container `id`, the `name`, the `image` and the current
`status`

```bash
89532c3f0369.container.	3600	IN	TXT	"id=89532c3f0369" "name=drunk_newton" "image=redis" "status=running"
```

### A records

They include the ips the container has.

```bash
89532c3f0369.container.	3600	IN	A	172.17.0.3
```
### SRV records

Here's where things get a little tricky. Trying to expose more information there
is a quirk in how the information is being delivered.

Let's use this example:

```bash
89532c3f0369.container.	3600	IN	SRV	0 0 6237 89532c3f0369.container.
89532c3f0369.container.	3600	IN	SRV	0 1 32768 localhost.localdomain.
89532c3f0369.container.	3600	IN	SRV	0 2 6379 89532c3f0369.container.
89532c3f0369.container.	3600	IN	SRV	0 3 0 localhost.localdomain.
```

In this scenario the container has two exposed ports `6237` and `6379`:

```bash
89532c3f0369.container.	3600	IN	SRV	0 0 6237 89532c3f0369.container.
89532c3f0369.container.	3600	IN	SRV	0 2 6379 89532c3f0369.container.
```

And we see that the host machine has two ports as well, `32768` and `0`:

```bash
89532c3f0369.container.	3600	IN	SRV	0 1 32768 localhost.localdomain.
89532c3f0369.container.	3600	IN	SRV	0 3 0 localhost.localdomain.
```

This actually means that the port `6237` in the container is exposed in the port
`32768` in the host machine and that the `6379` does not have an exposed port.
