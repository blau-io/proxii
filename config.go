package main

import (
	"flag"
)

type Config struct {
	etcdAddress  string
	etcdCaFile   string
	etcdCertFile string
	etcdKeyFile  string
	http2        bool
	port         int
	rootDir      string
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.etcdAddress, "etcdAddress", "http://127.0.0.1:4001/",
		"etcd address")
	flag.StringVar(&config.etcdCaFile, "etcdCaFile", "", "Path to etcd CA file")
	flag.StringVar(&config.etcdCertFile, "etcdCertFile", "",
		"Path to etcd Cert file")
	flag.StringVar(&config.etcdKeyFile, "etcdKeyFile", "",
		"Path to etcd Key file")
	flag.BoolVar(&config.http2, "http2", false, "Use HTTP2")
	flag.IntVar(&config.port, "port", 80, "Port to listen on")
	flag.StringVar(&config.rootDir, "rootDir", "/proxii",
		"Root directory of Proxii")

	flag.Parse()
	return config
}
