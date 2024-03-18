package driver

type Backend interface {
	Init()
	Dispose()
	Seed(documentNumber int)
	Search3(regex string) []Car
	CountDocuments() int64
	InsertOne(car Car) *Car
	GetByUUID(UUID string) *Car
}

type Car struct {
	UUID        string   `bson:"uuid,omitempty" json:"uuid"`
	PlateNumber string   `bson:"rendszam,omitempty" json:"rendszam" fake:"{regex:[A-Z]{7}}-{regex:[0-9]{1}}" form:"rendszam" binding:"required,min=7,max=7"`
	Owner       string   `bson:"tulajdonos,omitempty" json:"tulajdonos" fake:"{name}" form:"tulajdonos" binding:"required,alphanum,min=1,max=32"`
	ValidUntil  string   `bson:"forgalmi_ervenyes,omitempty" json:"forgalmi_ervenyes" fake:"{date}" format:"2006-01-02" form:"forgalmi_ervenyes" binding:"required,min=10,max=10"`
	Data        []string `bson:"adatok,omitempty" json:"adatok" fakesize:"3"`
}
