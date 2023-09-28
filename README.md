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

# How to deploy manually

- before commit/merge changes to `main`, bump `pkg/latest-version/latest-version.txt` version
- commit/merge changes to main
- create tag = `latest-version.txt`, e.g. `v0.3.0`
- `git push --tags`
- build binaries: `make build-n-compress-all` (they gitingored)
- make release for latest tag in github and attach `*.tar.gz` binaries 
- update `Dockerfile` in [nginx-reverse-proxy](https://github.com/oleksii-honchar/nginx-reverse-proxy) to fetch new `nrp-cli` version