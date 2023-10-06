package config

type NrpServiceConfig struct {
	SchemaVersion    string `yaml:"schemaVersion"`
	Name             string `yaml:"name"`
	ServiceIP        string `yaml:"serviceIp"`
	ServicePort      int    `yaml:"servicePort"`
	DomainName       string `yaml:"domainName"`
	DomainRegistrant string `yaml:"domainRegistrant,omitempty"`
	CORS             string `yaml:"cors"`
	BlockExploits    string `yaml:"blockExploits"`
	HTTPS            struct {
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
type NrpCronConfig struct {
	ConfigPath string `yaml:"configPath"`
}
type NrpPublicIpConfig struct {
	CheckAndUpdate string `yaml:"checkAndUpdate"`
	Schedule       string `yaml:"schedule"`
	DataPath       string `yaml:"dataPath"`
	DryRun         string `yaml:"dryRun"`
}

type NrpConfig struct {
	PublicIp    NrpPublicIpConfig   `yaml:"public-ip"`
	Supervisor  NrpSupervisorConfig `yaml:"supervisor"`
	Nginx       NrpNginxConfig      `yaml:"nginx"`
	Squid       NrpSquidConfig      `yaml:"squid"`
	Dnsmasq     NrpDnsmasqConfig    `yaml:"dnsmasq"`
	Letsencrypt LetsencryptConfig   `yaml:"letsencrypt"`
	Cron        NrpCronConfig       `yaml:"cron"`
	Services    []NrpServiceConfig  `yaml:"services"`
}

type NewCertRequest struct {
	DryRun   string
	BasePath string
	CertName string
	Email    string
	Domain   string
}
