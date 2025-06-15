package model

type SensorData struct {
	Temp struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
		Med float64 `json:"med"`
	} `json:"temp"`
	Pressure struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
		Med float64 `json:"med"`
	} `json:"pressure"`
	Humidity struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
		Med float64 `json:"med"`
	} `json:"humidity"`
}
