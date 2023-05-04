// All project structurs.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Vasily-van-Zaam/ushortener/pkg/shorter"
)

// User structure.
type UserData string

// User data constant.
const (
	USERDATA UserData = "user_data"
)

// Link structure. For save in DB and returns to client.
type Link struct {
	ID        int    `db:"id" json:"id"`
	Link      string `db:"link" json:"link"`
	ShortLink string `db:"short_link" json:"short_link"`
	UUID      string `db:"uuid" json:"uuid"`
	UserID    int    `db:"user_id" json:"user_id"`
	Deleted   bool   `db:"deleted" json:"deleted"`
}

// Function covert ID number to string 59.
func (l *Link) ConverID() string {
	sh := shorter.NewShorter59()
	id := sh.Convert(fmt.Sprint(l.ID))
	return id
}

// Error conflict struct.
type ErrConflict struct{}

// Function error.Error().
func (e *ErrConflict) Error() string {
	return "conflict"
}

// Create new error conflict.
func NewErrConflict() *ErrConflict {
	return &ErrConflict{}
}

// User struct. For autorization data.
type User struct {
	ID string `db:"id" json:"id"`
}

// Set id user from context.
func (u *User) SetUserIDFromContext(ctx context.Context) error {
	v, ok := ctx.Value(USERDATA).(User)
	if !ok {
		return errors.New("ERROR COOCIES")
	}
	u.ID = v.ID
	return nil
}

// Request API Shorten.
type RequestAPIShorten struct {
	URL string `json:"url"`
}

// Response API Shorten.
type ResponseAPIShorten struct {
	Result string `json:"result"`
}

// Response urls user.
type ResponseAPIUserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// Request  url Batch user.
type RequestAPIShortenBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Response url Batch user.
type ResponseAPIShortenBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Main config struct.
type Config struct {
	ServerAddress    string `json:"server_address" env:"SERVER_ADDRESS"`
	BaseURL          string `json:"base_url" env:"BASE_URL"`
	Filestore        string `json:"file_storage_path" env:"FILE_STORAGE_PATH"`
	SqliteDB         string `env:"SQLITE_DB"`
	ServerTimeout    int64  `env:"SERVER_TIMEOUT" envDefault:"100"`
	ExpiresDayCookie int64  `env:"EXPIRES_DAY_COOKIE" envDefault:"365"`
	SecretKey        string `env:"SECRET_KEY" envDefault:"secretkey"`
	DataBaseDNS      string `json:"database_dsn" env:"DATABASE_DSN"`
	EnableHTTPS      bool   `json:"enable_https" env:"ENABLE_HTTPS"`
	ConfigPath       string `env:"CONFIG"`
}

// Load config from file.
func (c *Config) UpdateFromJSON() {
	file, errFile := os.ReadFile(c.ConfigPath)
	if errFile != nil {
		log.Println(errFile)
	}
	err := json.Unmarshal(file, c)
	if err != nil {
		log.Println(errFile)
	}
}

// For list flags.
type flagVars struct {
	name      string
	value     string
	valueBool bool
	usage     string
	field     *string
	fieldBool *bool
}

// Set default values config.
func (c *Config) SetDefault() {
	flags := []flagVars{
		{name: "c", usage: "config path use as: -c ./config.js", field: &c.ConfigPath},
		{name: "config", usage: "use as: -config ./config.js", field: &c.ConfigPath},
		{name: "b", usage: "use as: -b http://example.com", value: "http://localhost:8080", field: &c.BaseURL},
		{name: "a", usage: "use as: -a 127.0.0.1:8080 or localhost:8080", value: "127.0.0.1:8080", field: &c.ServerAddress},
		{name: "f", usage: "use as: -f ./store.csv", field: &c.Filestore},
		{name: "d", usage: "path DB use as: -d ", field: &c.DataBaseDNS},
		{name: "s", usage: "enabled http, use as: -s on", fieldBool: &c.EnableHTTPS},
	}
	emptyVar := ""
	emptyBoolVar := false
	for _, f := range flags {
		switch {
		case f.field == nil:
			if !*f.fieldBool {
				flag.BoolVar(f.fieldBool, f.name, f.valueBool, f.usage)
			} else {
				flag.BoolVar(&emptyBoolVar, f.name, f.valueBool, f.usage)
			}
		case f.fieldBool == nil:
			if *f.field == "" {
				flag.StringVar(f.field, f.name, f.value, f.usage)
			} else {
				flag.StringVar(&emptyVar, f.name, f.value, f.usage)
			}
			if c.ConfigPath != "" && f.name == "c" || f.name == "config" {
				c.UpdateFromJSON()
			}
		}
	}

	flag.Parse()
}

// Logger response.
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

// Logger request.
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

// Struct for bufer deleete URL.
type BuferDeleteURL struct {
	IDS  []*string
	User *User
	Ctx  context.Context
}

// Un Converter string59 to in.
func (b *BuferDeleteURL) UnConvertIDS() []*string {
	sh := shorter.NewShorter59()
	ids := make([]*string, len(b.IDS))
	for i, id := range b.IDS {
		uid := sh.UnConvert(*id)
		ids[i] = &uid
	}
	return ids
}
