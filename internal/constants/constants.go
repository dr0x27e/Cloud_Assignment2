package constants

const (
	Currency_API  = "http://129.241.150.113:9090/currency/"
	Country_API   = "http://129.241.150.113:8080/v3.1/alpha/"
	OpenMeteo_API = "https://api.open-meteo.com/v1/forecast"

	Registration  = "Registrations"
	Webhooks      = "webhooks"
	SubCollection = "subCollection"
	Path          = "Path"

	WebhookEndpoint = "/registration"
	ServiceEndpoint = "/invocation"

	TEST_COUNTRY   = "http://129.241.150.113:8080/v3.1/alpha/no?fields=name"
	TEST_OPENMETEO = "https://api.open-meteo.com/v1/forecast?latitude=62&longitude=10&hourly=temperature_2m"
	TEST_CURRENCY  = "http://129.241.150.113:9090/currency/nok"
)

const (
	EventRegister = "REGISTER"
	EventChange   = "CHANGE"
	EventDelete   = "DELETE"
	EventInvoke   = "INVOKE"
)

var ValidEvents = map[string]bool{
	EventRegister: true,
	EventChange:   true,
	EventDelete:   true,
	EventInvoke:   true,
}

var AllEvents = []string{
	EventRegister,
	EventChange,
	EventDelete,
	EventInvoke,
}
