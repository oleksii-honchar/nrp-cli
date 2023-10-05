* [Basic usage](#basic-usage)
* [CMD options](#cmd-options)
* [Adding new service to nrp.yaml](#adding-new-service-to-nrpyaml)
    + [HTTPS](#https)
* [How to deploy manually](#how-to-deploy-manually)
* [Solution Design](#solution-design)
* [Troubleshooting](#troubleshooting)
# nrp-cli
[Nginx-Reverse-Proxy](https://github.com/oleksii-honchar/nginx-reverse-proxy) cli tool for creating `nginx` proxy-path configs with simple yaml configuration and HTTPS support.

## Basic usage

- Create empty `nrp.yaml` file and run `make run-all`. Go to `localhost` in browser. You should see `nginx-more` default page.

- Assume your services executed on the same host as nginx, and host has local IP = `192.168.0.12`. Also you had configured `service1.domain.tld` & `service1.domain.tld` redirect to your ISP public IP.

- Create `nrp.yaml` file in root repo folder:
```yaml
schemaVersion: 0.3.0
letsencrypt:
  email: "you-name@gmail.com"
  dryRun: false
services:
- name: service1
  serviceIp: 192.168.0.12
  servicePort: 9000
  domainName: service1.domain.tld
  cors: true
- name: service2
  serviceIp: 192.168.0.12
  servicePort: 9100
  domainName: service2.domain.tld
  cors: true
  https: 
    use: true
    force: true 
    hsts: true
```

- If you don't have your ISP Nat-loopback enabled, then add `127.0.0.1 service2.domain.tld ` to `/etc/hosts` 

```bash
make run-all
```
- `nrp-cli` will requests certificates for `service2` and then start nginx.
- If you are lucky, you will be able to request both of your services via domain name
- In case you don't have your ISP Nat-loopback enabled, you can reach them only by adding to `/etc/hosts`. This case covered in [nginx-reverse-proxy](https://github.com/oleksii-honchar/nginx-reverse-proxy) project as well as more complex services setup.
- It is not recommeneded to change `defaults.prod` values as they directly affect cli behaviour
- You can use `letsencrypt.dryRun = trye` option to verify that `certbot` is able to perform domain verification

## CMD options

- `-h, -help` - shows cmd help
- `-v, -version` - show build version
- `-log-level` - info(default)|error|warn|debug
- `-config` - path to `nrp.yaml` - './nrp.yaml'(default)
- `-defaults-mode` - dev|prod(default) with nginx & letsencrypt param wil be used. "Dev" used only for development
- `-certbot-wait` - debug only option, will make `nrp-cli` sleep for 5 min right before making certbot request

## Adding new service to nrp.yaml

Just add new array item in `nrp.yaml`:
- Keep `name` unique 
- `domainName` should contain only single domain, multiple domains not tested
- rest of the options will add related nginx includes to the service config

```
- name: service2
  serviceIp: 192.168.0.12
  servicePort: 9100
  domainName: service2.domain.tld
  cors: true
  https: 
    use: true
    force: true 
    hsts: true
```

## How to deploy manually

- before commit/merge changes to `main`, bump `pkg/latest-version/latest-version.txt` version
- commit/merge changes to main
- create tag = `latest-version.txt`, e.g. `v0.3.0`
- `git push --tags`
- build binaries: `make build-n-compress-all` (they gitingored)
- make release for latest tag in github and attach `*.tar.gz` binaries 
- update `Dockerfile` in [nginx-reverse-proxy](https://github.com/oleksii-honchar/nginx-reverse-proxy) to fetch new `nrp-cli` version

## Solution Design
Here is the nginx configuration decomposition in chunks from which then every service config composed:
![nginx-config-structure](./docs/nrp-nginx-config-structure.jpg)

Here is the flow diagram for main logic:
![nrp-flow-diagram](./docs/nrp-flow-diagram.jpg)

## Troubleshooting

- When used without NRP, `nrp-cli` will fail when requesting SSL cert via Certbot. Domains need to be added to `/etc/hosts`.
- Explore `Makefile` targets to help you debug nginx & nrp-cli behaviour