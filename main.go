package main

import (
	"flag"

	"whoamiv2/pkg/dns"
)

func main() {
	serv := dns.NewServer()
	domain := flag.String("d", "domain.kr", "-d=mydomain.com")

	flag.Parse()

	serv.Start(*domain)
}
