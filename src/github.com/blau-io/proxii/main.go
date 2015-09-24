package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type proxii struct {
	config *Config
	etcd   *etcdConnector
}

func (p *proxii) handler(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")[0]

	uri, err := p.etcd.resolve(host)
	if err != nil {
		log.Println("Error while looking up host: ", err)
		return
	}

	proxy := newReverseProxy(uri)
	proxy.ServeHTTP(w, r)
}

func main() {
	p := newProxii(parseFlags())
	http.HandleFunc("/", p.handler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", p.config.port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func newProxii(config *Config) *proxii {
	etcd, err := newEtcdConnector(config)
	if err != nil {
		panic(err)
	}

	etcd.init()
	return &proxii{config, etcd}
}