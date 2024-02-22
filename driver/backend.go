package driver

type Backend interface {
	Init()
	Dispose()
	Seed(documentNumber int)
	Search3(regex string)
}

type Car struct {
	UUID        string   `bson:"uuid,omitempty" json:"uuid"`
	PlateNumber string   `bson:"rendszam,omitempty" json:"rendszam" fake:"{regex:[A-Z]{7}}-{regex:[0-9]{1}}"`
	Owner       string   `bson:"tulajdonos,omitempty" json:"tulajdonos" fake:"{name}"`
	ValidUntil  string   `bson:"forgalmi_ervenyes,omitempty" json:"forgalmi_ervenyes" fake:"{date}" format:"2006-01-02"`
	Data        []string `bson:"adatok,omitempty" json:"adatok" fakesize:"3"`
}
