package models

import "time"

type ErrorResponse struct {
	StatusCode   int    `json:"statusCode"`
	ErrorMessage string `json:"errorMessage"`
}

type Article struct {
	Id                string    `json:"id"`
	Title             string    `json:"title"`
	SubTitle          string    `json:"subTitle"`
	Content           string    `json:"content"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}
