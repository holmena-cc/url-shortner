package models

import "time"

type Visit struct {
    ClickID   int       `db:"click_id" json:"clickId"`
    UrlID     int       `db:"url_id" json:"urlId"` // foreign key to Url
    ClickDate time.Time `db:"click_date" json:"clickDate"`
    Referrer  string    `db:"referrer" json:"referrer"`
    IPAddress string    `db:"ip_address" json:"ipAddress"`
    Country   string    `db:"country" json:"country"`
}