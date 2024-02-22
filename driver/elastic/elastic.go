package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/cyrip/monGO/driver"
	"github.com/google/uuid"
	elastic "github.com/olivere/elastic/v7"
)

const ELASTIC_INDEX_NAME string = "cars"
const ELASTIC_URL = "http://127.0.0.1:9200"

var client *elastic.Client

type Elastic struct {
	indexName     string
	elasticClient *elastic.Client
}

func (this *Elastic) Init(indexName string) *elastic.Client {
	if this.elasticClient == nil {
		var err error
		this.elasticClient, err = elastic.NewClient(
			elastic.SetURL(ELASTIC_URL),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Fatal(err)
		}
		this.indexName = indexName
	}
	return this.elasticClient
}

func (this *Elastic) Dispose() {
	//return
}

func (this *Elastic) CreateIndex() {
	mapping := `{
		"settings": {
			"number_of_shards": 2,
			"number_of_replicas": 1
		},
		"mappings": {
			"properties": {
				"rendszam": { "type": "keyword" },
				"tulajdonos": { "type": "keyword" },
				"forgalmi_ervenyes": { "type": "date" },
				"adatok": { "type": "keyword" }
			}
		}
	}`

	// Create an index with the defined settings and mappings
	createIndex, err := this.elasticClient.CreateIndex(this.indexName).BodyString(mapping).Do(context.Background())
	if err != nil {
		log.Fatalf("Failed to create index: %s", err)
	}
	if !createIndex.Acknowledged {
		log.Fatal("Create index not acknowledged")
	}

	log.Println("Index created successfully")
}

func (this *Elastic) DeleteIndex() {
	deleteIndex, err := this.elasticClient.DeleteIndex(this.indexName).Do(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
		log.Println("Delete index not acknowledged")
	} else {
		fmt.Println("Index deleted")
	}
}

func (this *Elastic) GetAllDocuments() {
	// Initialize scrolling over documents
	scroll := this.elasticClient.Scroll(this.indexName).Size(100) // Adjust size as needed
	for {
		results, err := scroll.Do(context.Background())
		if err == io.EOF {
			log.Println("All documents retrieved")
			break
		}
		if err != nil {
			log.Fatalf("Error retrieving documents: %s", err)
		}

		// Iterate through results
		for _, hit := range results.Hits.Hits {
			var doc map[string]interface{}
			if err := json.Unmarshal(hit.Source, &doc); err != nil {
				log.Fatalf("Error deserializing document: %s", err)
			}
			// Process your document (doc) here
			log.Printf("Doc ID: %s, Doc: %+v\n", hit.Id, doc)
		}
	}
}

func (this *Elastic) Search3(term string) {
	query := elastic.NewBoolQuery().Should(
		elastic.NewRegexpQuery("rendszam", term),
		elastic.NewRegexpQuery("tulajdonos", term),
		elastic.NewRegexpQuery("adatok", term),
	)

	searchResult, err := this.elasticClient.Search().
		Index(this.indexName).
		Query(query).
		Pretty(true).
		Size(100).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	fmt.Printf("Found %d documents\n", searchResult.TotalHits())

	for _, hit := range searchResult.Hits.Hits {
		//var doc map[string]interface{}
		var doc driver.Car
		err := json.Unmarshal(hit.Source, &doc)
		if err != nil {
			log.Fatalf("Error deserializing hit to document: %s", err)
		}
		doc.UUID = hit.Id
		fmt.Printf("Document ID: %s, Fields: %+v\n", hit.Id, doc)
	}
}

func (this *Elastic) GetUUID(name string) string {
	namespaceDNS := uuid.NameSpaceDNS
	uuidV5 := uuid.NewSHA1(namespaceDNS, []byte(name))
	return uuidV5.String()
}

func (this *Elastic) AddDocument(doc driver.Car) bool {
	uuid5 := this.GetUUID(doc.PlateNumber)
	indexResponse, err := this.elasticClient.Index().
		Index(this.indexName).
		BodyJson(doc).
		Id(uuid5).
		Do(context.Background())
	if err != nil {
		log.Println(err)
		return false
	}

	log.Printf("Indexed document %s to index %s\n", indexResponse.Id, indexResponse.Index)
	return true
}

func (this *Elastic) Test1(plateNumber string) {
	car := driver.Car{
		//UUID5:       uuid5,
		PlateNumber: plateNumber,
		Owner:       "KZ",
		ValidUntil:  "2024-01-01",
		Data:        []string{"data1", "data2", "data3"},
	}
	this.AddDocument(car)
	this.Search3(".*ABC.*")
}

func (this *Elastic) InsertOne(car driver.Car) {
	this.AddDocument(car)
}

func (this *Elastic) Seed(count int) {
	gofakeit.Seed(rand.Intn(10000))
	inserted := 0
	notInserted := 0
	for i := 0; i < count; i++ {
		car := this.getFakeCar()
		log.Println(car)
		if this.AddDocument(car) {
			inserted = inserted + 1
		} else {
			notInserted = notInserted + 1
		}
	}
	log.Printf("Inserted/Not inserted", inserted, notInserted)
}

func (this *Elastic) getFakeCar() driver.Car {
	var fakeCar driver.Car
	err := gofakeit.Struct(&fakeCar)
	if err != nil {
		log.Fatal(err)
	}

	fakeCar.ValidUntil = fakeCar.ValidUntil[0:10]
	fakeCar.Data[0] = gofakeit.Color()
	fakeCar.Data[1] = gofakeit.City()
	fakeCar.Data[2] = gofakeit.BeerName()

	return fakeCar
}

func (this *Elastic) CountDocuments() int {
	count, err := this.elasticClient.Count().
		Index(this.indexName).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting document count: %s", err)
	}

	fmt.Printf("Document count in '%s': %d\n", this.indexName, count)
	return int(count)
}

func (this *Elastic) DeleteDocument(docID string) {
	deleteResponse, err := this.elasticClient.Delete().
		Index(this.indexName).
		Id(docID).
		Do(context.Background())

	if err != nil {
		log.Printf("Error deleting document: %s", err.Error())
	}
	fmt.Printf("Document ID %s deleted, result: %s\n", docID, deleteResponse.Result)
}
