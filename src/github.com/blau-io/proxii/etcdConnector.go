package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

type etcdConnector struct {
	etcd    client.KeysAPI
	hosts   map[string]string
	rootDir string
}

func newEtcdConnector(c *Config) (*etcdConnector, error) {
	var transport client.CancelableTransport
	if c.etcdCertFile == "" && c.etcdKeyFile == "" {
		transport = client.DefaultTransport
	} else {
		tlsCert, err := tls.LoadX509KeyPair(c.etcdCertFile, c.etcdKeyFile)
		if err != nil {
			return nil, err
		}

		certBytes, err := ioutil.ReadFile(c.etcdCaFile)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(certBytes)

		transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{tlsCert},
				RootCAs:      caCertPool,
			},
			TLSHandshakeTimeout: 10 * time.Second,
		}
	}

	etcdClient, err := client.New(client.Config{
		Endpoints:               []string{c.etcdAddress},
		Transport:               transport,
		HeaderTimeoutPerRequest: time.Second,
	})

	if err != nil {
		return nil, err
	}

	etcd := client.NewKeysAPI(etcdClient)
	hosts := make(map[string]string)
	rootDir := c.rootDir

	return &etcdConnector{etcd, hosts, rootDir}, nil
}

func (con *etcdConnector) init() {
	resp, err := con.etcd.Get(context.Background(), con.rootDir,
		&client.GetOptions{true, false, true})

	if err == nil {
		for _, node := range resp.Node.Nodes {
			con.update(node, resp.Action)
		}
	}

	log.Printf("Watching %s", con.rootDir)
	go con.watch()
}

func (con *etcdConnector) resolve(host string) (*url.URL, error) {
	host = path.Join(con.rootDir, host)

	value := con.hosts[host]
	if value == "" {
		return nil, errors.New(fmt.Sprintf("No entry found for %s", host))
	}

	uri, err := url.Parse(value)
	if err != nil {
		return nil, err
	}

	return uri, nil
}

func (con *etcdConnector) update(node *client.Node, action string) {
	if action == "delete" || action == "expire" {
		log.Printf("Deleted Host '%v'", node.Key)
		delete(con.hosts, node.Key)
		return
	}

	if node.Value != con.hosts[node.Key] {
		con.hosts[node.Key] = node.Value
		log.Printf("Found Host  '%v' with Value '%v'", node.Key, node.Value)
	}
}

func (con *etcdConnector) watch() {
	for {
		resp, err := con.etcd.Watcher(con.rootDir, &client.WatcherOptions{
			(uint64)(0), true}).Next(context.Background())

		if err != nil {
			log.Printf("Error while watching /proxii/: %v", err)
			time.Sleep(time.Second)
		} else {
			con.update(resp.Node, resp.Action)
		}
	}
}
