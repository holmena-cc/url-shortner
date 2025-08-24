package models

import "time"

type Url struct {
    UrlID        int       `db:"url_id" json:"urlId"`
    OriginalUrl  string    `db:"original_url" json:"originalUrl"`
    ShortCode    string    `db:"short_code" json:"shortCode"`
    CustomAlias  *string   `db:"custom_alias" json:"customAlias,omitempty"` // optional
    CreationDate time.Time `db:"creation_date" json:"creationDate"`
    UserID       int       `db:"user_id" json:"userId"` // foreign key to User
}