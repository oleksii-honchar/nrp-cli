module npr/npr-cli

go 1.21.1

require (
	beaver/blablo v0.0.0-00010101000000-000000000000
	nrp/config-processor v0.0.0-00010101000000-000000000000
)

require gopkg.in/yaml.v2 v2.4.0 // indirect

replace beaver/blablo => ./blablo

replace nrp/config-processor => ./pkg/config-processor
