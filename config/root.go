package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Twitter struct {
	Ticker            int64
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	BearerToken       string

	StartTime int64
}

type Config struct {
	Docker struct {
		Init    bool
		Targets []string // []string{"es01", "kibana"}
	}

	Server struct {
		Port string
	}

	Repository struct {
		ElasticSearch struct {
			URI      string
			User     string
			Password string
			ApiKey   string
		}
	}

	Cache struct {
		Redis struct {
			DataSource string
			DB         int
			Password   string
			UserName   string
			Beta       float64
		}

		Local struct {
			LocalCacheTTL int64
		}
	}

	Twitter map[string]Twitter // key is search keyword

	OllaMa struct {
		Model string
		Url   string
	}
}

func NewConfig(path string) Config {
	c := new(Config)

	if file, err := os.Open(path); err != nil {
		panic(err)
	} else {
		defer file.Close()
		if err = toml.NewDecoder(file).Decode(c); err != nil {
			panic(err)
		} else {
			return *c
		}
	}
}
