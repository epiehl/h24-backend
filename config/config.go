package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/spf13/viper"
)

type config struct {
	Database struct {
		User     string
		Password string
		DBName   string
		Host     string
		Port     string
		Sslmode  string
		Timezone string
	}
	Server struct {
		Host       string
		Port       string
		Scheme     string
		Production bool
	}
	Cognito struct {
		PoolID string
	}
	AWS struct {
		Region string
		SES    struct {
			From string
		}
	}
	Auth struct {
		JWksUrl string
	}
	Aggregator struct {
		Outlet struct {
			Endpoint struct {
				Scheme   string
				Host     string
				Location string
			}
		}
	}
	H24Connector struct {
		Endpoint string
	}
}

var C config

func ReadConfig() {
	Config := &C

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join("$GOPATH", "src", "github.com", "epiehl93", "h24-notifier", "config"))
	viper.SetEnvPrefix("H24")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println(err)
		} else {
			log.Fatalln(err)
		}
	}

	if err := viper.Unmarshal(&Config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	spew.Dump(1)
}
