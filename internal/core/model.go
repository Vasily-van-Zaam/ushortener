package core

import "flag"

const MAINDOMAIN = "http://localhost:8080/"

type Link struct {
	ID        int    `db:"id" json:"id"`
	Link      string `db:"link" json:"link"`
	ShortLink string `db:"short_link" json:"short_link"`
}

type RequestAPIShorten struct {
	URL string `json:"url"`
}

type ResponseAPIShorten struct {
	Result string `json:"result"`
}

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
	Filestore     string `env:"FILE_STORAGE_PATH"`
	SqliteDB      string `env:"SQLITE_DB"`
}

func (c *Config) SetDefault() {

	if c.BaseURL == "" {
		flag.StringVar(&c.BaseURL, "b", "http:/localhost:8080", "use as http://example.com")
	} else {
		v := ""
		flag.StringVar(&v, "b", "http:/localhost:8080", "use as http://example.com")
	}

	if c.ServerAddress == "" {
		flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "use as 127.0.0.1:8080 or localhost:8080")
	} else {
		v := ""
		flag.StringVar(&v, "a", "localhost:8080", "use as 127.0.0.1:8080 or localhost:8080")
	}

	if c.Filestore == "" {
		flag.StringVar(&c.Filestore, "f", "./store", "path to file ./store.csv")
	} else {
		v := ""
		flag.StringVar(&v, "f", "./store", "path to file ./store.csv or other")
	}

	flag.Parse()

}
