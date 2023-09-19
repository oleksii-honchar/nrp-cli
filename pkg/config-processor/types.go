package configProcessor

type NrpServiceConfig struct {
	Name          string `yaml:"name"`
	ServiceIP     string `yaml:"serviceIp"`
	ServicePort   int    `yaml:"servicePort"`
	DomainName    string `yaml:"domainName"`
	CORS          bool   `yaml:"cors,omitempty"`
	BlockExploits bool   `yaml:"blockExploits,omitempty"`
	HTTPS         struct {
		Use   bool `yaml:"use,omitempty"`
		Force bool `yaml:"force,omitempty"`
		HSTS  bool `yaml:"hsts,omitempty"`
	} `yaml:"https,omitempty"`
}

type NrpNginxConfig struct {
	StartCmd   string `yaml:"startCmd"`
	StopCmd    string `yaml:"stopCmd"`
	StatusCmd  string `yaml:"statusCmd"`
	LogsCmd    string `yaml:"logsCmd"`
	ConfigPath string `yaml:"configPath"`
}

type LetsencryptConfig struct {
	DryRun             bool   `yaml:"dryRun"`
	Email              string `yaml:"email"`
	BasePath           string `yaml:"basePath"`
	CertFilesPath      string `yaml:"certFilesPath"`
	AcmeChallengePath  string `yaml:"acmeChallengePath"`
	RequestCertCmdTmpl string `yaml:"requestCertCmdTmpl"`
	RenewCertCmd       string `yaml:"renewCertCmd"`
}

type NrpConfig struct {
	Nginx       NrpNginxConfig     `yaml:"nginx"`
	Letsencrypt LetsencryptConfig  `yaml:"letsencrypt"`
	Services    []NrpServiceConfig `yaml:"services"`
}

type NewCertRequest struct {
	DryRun   bool
	BaseDir  string
	CertName string
	Email    string
	Domain   string
}
