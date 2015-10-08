#Proxii
Proxii is a small reverse Proxy which stores its configuration in etcd. It tries to mimic `proxy_pass` from Nginx.

##Disclaimer:
I decided to replace this in production with Nginx + Confd. I'm leaving the code here so people can observe it if they want to. I don't plan on updating this any further.

This project is heavily inspired by [Vulcand](https://vulcand.io/) and [Gogeta](https://github.com/arkenio/gogeta). However, both projects failed to replace a simple Nginx `proxy_pass` configuration, thus creating the need for a better solution.

The main difference between Gogeta and this is that this proxy will follow redirects.
