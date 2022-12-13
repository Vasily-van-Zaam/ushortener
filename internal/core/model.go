package core

const MAINDOMAIN = "http://localhost:8080/"

type Link struct {
	ID        int    `db:"id" json:"id"`
	Link      string `db:"link" json:"link"`
	ShortLink string `db:"short_link" json:"short_link"`
}

type RequestApiShorten struct {
	URL string `json:"url"`
}

type ResponseApiShorten struct {
	Result string `json:"result"`
}
