package models

type Permission struct {
	BaseModel
	UserID    string
	Name      string `json:"name"`
	CreatedBy string
}
