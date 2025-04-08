package structs

import (
	"cloud.google.com/go/firestore"
)

type StatusResponse struct {
	CountriesAPI 	     int
	MeteoAPI 			 int
	CurrencyAPI 	     int
	PythonAPI            int
	WebhooksDatabase     int
	RegistrationDatabase int
	Webhooks             int
	Version              string
	Uptime               int64
}

type WebhookRegistrationModel struct {
	Url     string `json:"url"`
	Country string `json:"country"`
	Event   string `json:"event"`
}

type Configuration struct {
	Country    string   `firestore:"country" json:"country"`
	ISOCode    string   `firestore:"isoCode" json:"isoCode"`
	Features   Features `firestore:"features" json:"features"`
	LastChange string   `firestore:"lastChange" json:"lastChange"`
}

type Features struct {
	Temperature      bool     `firestore:"temperature"  json:"temperature"`
	Precipitation    bool     `firestore:"precipitation" json:"precipitation"`
	Capital          bool     `firestore:"capital" json:"capital"`
	Coordinates      bool     `firestore:"coordiantes" json:"coordinates"`
	Population       bool     `firestore:"population" json:"population"`
	Area             bool     `firestore:"area" json:"area"`
	TargetCurrencies []string `firestore:"targetCurrencies" json:"targetCurrencies"`
}

type Country struct {
	Capital    []string  `json:"capital"`
	Latitude   []float64 `json:"latlng"`
	Population int       `json:"population"`
	Area       float64   `json:"area"`
	Currency   map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

type FirebaseClient struct {
	Client *firestore.Client
}
