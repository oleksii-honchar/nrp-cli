# nrp-cli
Nginx-Reverse-Proxy cli tool

# CMD options
`-log-level`
`-config`
# Adding new service to nrp.yaml

## HTTPS
If https needed, first you need to add service domain forward to `127.0.0.1` in `/etc/hosts`. Otherwise certbot will be unable to retreive certificates. This happens, when your ISP connection doesn't support NAT-loopback. For example:
```text
# /etc/hosts

127.0.0.1 your.domain.tld
```