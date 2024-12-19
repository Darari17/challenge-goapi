package models

type Customer struct {
	Id          string `json:"id"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

type CustomerUpdate struct {
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phoneNumber"`
	Address     *string `json:"address"`
}