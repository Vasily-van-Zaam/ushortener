package core

const MAINDOMAIN = "https://some-domain.com/"

type Link struct {
	Id        int    `db:"id" json:"id"`
	Link      string `db:"link" json:"link"`
	ShortLink string `db:"short_link" json:"short_link"`
}
