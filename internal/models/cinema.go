package models

type Cinema struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Location string `json:"location"`
}
