package configProcessor

type NrpServiceConfig struct {
	SchemaVersion string `yaml:"schemaVersion"`
	Name          string `yaml:"name"`
	ServiceIP     string `yaml:"serviceIp"`
	ServicePort   int    `yaml:"servicePort"`
	DomainName    string `yaml:"domainName"`
	CORS          string `yaml:"cors"`
	BlockExploits string `yaml:"blockExploits"`
	HTTPS         struct {
		Use   string `yaml:"use"`
		Force string `yaml:"force"`
		HSTS  string `yaml:"hsts"`
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
	DryRun             string `yaml:"dryRun"`
	Email              string `yaml:"email"`
	BasePath           string `yaml:"basePath"`
	CertFilesPath      string `yaml:"certFilesPath"`
	RequestCertCmdTmpl string `yaml:"requestCertCmdTmpl"`
	RenewCertCmd       string `yaml:"renewCertCmd"`
}

type NrpSquidConfig struct {
	ConfigPath string `yaml:"configPath"`
	Use        string `yaml:"use"`
	UseDnsmasq string `yaml:"useDnsmasq"`
	Port       int    `yaml:"port"`
}

type NrpDnsmasqConfig struct {
	ConfigPath string `yaml:"configPath"`
}
type NrpSupervisorConfig struct {
	ConfigPath string `yaml:"configPath"`
}

type NrpConfig struct {
	Nginx       NrpNginxConfig      `yaml:"nginx"`
	Letsencrypt LetsencryptConfig   `yaml:"letsencrypt"`
	Services    []NrpServiceConfig  `yaml:"services"`
	Squid       NrpSquidConfig      `yaml:"squid"`
	Dnsmasq     NrpDnsmasqConfig    `yaml:"dnsmasq"`
	Supervisor  NrpSupervisorConfig `yaml:"supervisor"`
}

type NewCertRequest struct {
	DryRun   string
	BasePath string
	CertName string
	Email    string
	Domain   string
}
