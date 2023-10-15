package configDefaults

import _ "embed"

var DefaultsDevMode = "dev"
var DefaultsProdMode = "prod"

//go:embed nrp.defaults.dev.yaml
var NrpConfigDevDefaults []byte

//go:embed nrp.defaults.prod.yaml
var NrpConfigProdDefaults []byte

//go:embed nrp-service.defaults.yaml
var NrpSvcConfigDefaults []byte
