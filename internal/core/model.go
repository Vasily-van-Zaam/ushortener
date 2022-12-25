package core

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"reflect"
)

type UserData string

const (
	USERDATA UserData = "user_data"
)

type Link struct {
	ID        int    `db:"id" json:"id"`
	Link      string `db:"link" json:"link"`
	ShortLink string `db:"short_link" json:"short_link"`
	UserID    string `db:"user_id" json:"user_id"`
}

type User struct {
	ID string `db:"id" json:"id"`
}

func (u *User) FromAny(v any) {
	if reflect.TypeOf(v) == reflect.TypeOf(User{}) {
		u.ID = v.(User).ID
	}
}

type RequestAPIShorten struct {
	URL string `json:"url"`
}

type ResponseAPIShorten struct {
	Result string `json:"result"`
}

type ResponseAPIUserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Config struct {
	ServerAddress    string `env:"SERVER_ADDRESS"`
	BaseURL          string `env:"BASE_URL"`
	Filestore        string `env:"FILE_STORAGE_PATH"`
	SqliteDB         string `env:"SQLITE_DB"`
	ServerTimeout    int64  `env:"SERVER_TIMEOUT" envDefault:"100"`
	ExpiresDayCookie int64  `env:"EXPIRES_DAY_COOKIE" envDefault:"365"`
	SecretKey        string `env:"SECRET_KEY" envDefault:"secretkey"`
}

func (c *Config) SetDefault() {
	emptyVar := ""
	if c.BaseURL == "" {
		flag.StringVar(&c.BaseURL, "b", "http://localhost:8080", "use as http://example.com")
	} else {
		flag.StringVar(&emptyVar, "b", "http://localhost:8080", "use as http://example.com")
	}
	if c.ServerAddress == "" {
		flag.StringVar(&c.ServerAddress, "a", "127.0.0.1:8080", "use as 127.0.0.1:8080 or localhost:8080")
	} else {
		flag.StringVar(&emptyVar, "a", "127.0.0.1:8080", "use as 127.0.0.1:8080 or localhost:8080")
	}
	if c.Filestore == "" {
		flag.StringVar(&c.Filestore, "f", "./store", "path to file ./store.csv")
	} else {
		flag.StringVar(&emptyVar, "f", "./store", "path to file ./store.csv or other")
	}

	flag.Parse()
}

func (c *Config) LogResponse(w http.ResponseWriter, r *http.Request, body any, status int) {
	configByte, _ := json.Marshal(&c)
	log.Print(
		"\n# START LOG RESPONSE #", "\n",
		"Accept-Encoding: ", r.Header.Get("Accept-Encoding"), "\n",
		"Content-Encoding: ", r.Header.Get("Content-Encoding"), "\n",
		"Content-Type: ", "application/json", "\n",
		"STATUS: ", status, "\n",
		"METHOD: ", r.Method, "\n",
		"PROTO: ", r.Proto, "\n",
		"URL: ", r.Host, r.URL.Path, "\n",
		"BODY: ", body, "\n",
		"CONFIG: ", string(configByte), "\n",
		"# END LOG RESPONSE #", "\n",
	)
}
func (c *Config) LogRequest(w http.ResponseWriter, r *http.Request, body any) {
	configByte, _ := json.Marshal(&c)
	log.Print(
		"\n# START LOG REQUEST #", "\n",
		"Accept-Encoding: ", r.Header.Get("Accept-Encoding"), "\n",
		"Content-Encoding: ", r.Header.Get("Content-Encoding"), "\n",
		"METHOD: ", r.Method, "\n",
		"PROTO: ", r.Proto, "\n",
		"URL: ", r.Host, r.URL.Path, "\n",
		"BODY: ", body, "\n",
		"CONFIG: ", string(configByte), "\n",
		"# END LOG REQUEST #", "\n",
	)
}
