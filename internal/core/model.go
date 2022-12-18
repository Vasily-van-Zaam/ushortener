package core

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
		c.BaseURL = "http://localhost:8080"
	}
	if c.ServerAddress == "" {
		c.ServerAddress = "http://localhost:8080"
	}
}
