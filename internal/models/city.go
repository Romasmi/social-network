package models

type City struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CountyCode    string `json:"countyCode"`
	StateProvince string `json:"stateProvince"`
}
