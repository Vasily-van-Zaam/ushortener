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
