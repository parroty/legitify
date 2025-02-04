package cmd

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type args struct {
	Token         string
	Endpoint      string
	Organizations []string
	Repositories  []string
	PoliciesPath  []string
	Namespaces    []string
	ColorWhen     string
	OutputFile    string
	ErrorFile     string
	OutputFormat  string
	OutputScheme  string
	ScorecardWhen string
	FailedOnly    bool
}

const (
	ArgErrorFile  = "error-file"
	ArgOutputFile = "output-file"
	ArgToken      = "github-token"
	ArgServerUrl  = "server-url"
)

const (
	EnvToken     = "github_token"
	EnvServerUrl = "server_url"
)

func (a *args) ApplyEnvVars() {
	if a.Token == "" {
		a.Token = viper.GetString(EnvToken)
	}

	if a.Endpoint == "" {
		a.Endpoint = viper.GetString(EnvServerUrl)
	}
}

func (a *args) AddCommonOptions(flags *pflag.FlagSet) {
	flags.StringVarP(&a.Token, ArgToken, "t", "", "token to authenticate with github (required unless environment variable GITHUB_TOKEN is set)")
	flags.StringVarP(&a.Endpoint, ArgServerUrl, "", "", "github endpoint to use instead of GitHub Cloud (can be set via the environment variable SERVER_URL)")
	flags.StringVarP(&a.OutputFile, ArgOutputFile, "o", "", "output file, defaults to stdout")
	flags.StringVarP(&a.ErrorFile, ArgErrorFile, "e", "error.log", "error log path")
}
