package banner

import "time"

type Banner struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	ImageUrl  string    `json:"image_url"`
	Index     int       `json:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
