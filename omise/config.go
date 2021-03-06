package main

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var config = &struct {
	viper *viper.Viper

	PKey string `mapstructure:"pkey"`
	SKey string `mapstructure:"skey"`

	APIEndpoint   string `mapstructure:"api-endpoint"`
	VaultEndpoint string `mapstructure:"vault-endpoint"`

	JSON         bool `mapstructure:"json"`
	IndentedJSON bool `mapstructure:"indented-json"`
}{}

var envMap = map[string]string{
	"pkey": "PUBKEY",
	"skey": "KEY",
}

func configure(flags *flag.FlagSet) {
	flags.String("pkey", "", "Omise API public key, defaults to $OMISE_PUBKEY.")
	flags.String("skey", "", "Omise API secret key, defaults to $OMISE_KEY.")
	flags.String("api-endpoint", "", "Omise API endpoint, defaults to api.omise.co or $OMISE_API_ENDPOINT")
	flags.String("vault-endpoint", "", "Omise Vault endpoint, defaults to vault.omise.co or $OMISE_VAULT_ENDPOINT")
	flags.Bool("json", false, "Output result as JSON.")
	flags.Bool("indented-json", false, "Output result as indented JSON.")
}

func bindViper(cmd *cobra.Command) error {
	v := viper.New()
	collectFlags(v, cmd)
	for _, key := range v.AllKeys() {
		if envkey, ok := envMap[key]; ok {
			v.BindEnv(key, "OMISE_"+envkey)
		} else {
			v.BindEnv(key, "OMISE_"+strings.Replace(strings.ToUpper(key), "-", "_", -1))
		}
	}

	configPath := os.Getenv("HOME") + "/.omise"
	if configFile, e := os.Open(configPath); e == nil {
		defer configFile.Close()

		v.SetConfigType("json")
		if e := v.ReadConfig(configFile); e != nil {
			return e
		}
	}

	if e := v.Unmarshal(config); e != nil {
		return e
	}

	config.viper = v
	return nil
}

func collectFlags(v *viper.Viper, cmd *cobra.Command) {
	v.BindPFlags(cmd.PersistentFlags())
	v.BindPFlags(cmd.Flags())
	for _, cmd := range cmd.Commands() {
		collectFlags(v, cmd)
	}
}
