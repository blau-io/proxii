#Proxii
Proxii is a small reverse Proxy which stores its configuration in etcd. It tries to mimic `proxy_pass` from Nginx.

This project is heavily inspired by [Vulcand](https://vulcand.io/) and [Gogeta](https://github.com/arkenio/gogeta). However, both projects failed to replace a simple Nginx `proxy_pass` configuration, thus creating the need for a better solution.

**This project is still in heavy development**

###Roadmap:
 - Support encrypted etcd nodes
 - Add unit tests
 - Add documentation
 - Support TLS
 - Support HTTP2
 - allow for more configuration


