# Inspector Gadget

_The Docker inspector_

![](http://clipset.20minutos.es/wp-content/uploads/2013/10/gadget.jpg)

## DNS Querying

```bash
dig [container_id].container @[DNS Address]
```

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
