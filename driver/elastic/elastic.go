package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/cyrip/monGO/config"
	"github.com/cyrip/monGO/driver"
	"github.com/cyrip/monGO/utils"
	elastic "github.com/olivere/elastic/v7"
)

const ELASTIC_INDEX_NAME string = "cars"
const ELASTIC_URL = "http://es00:9200"

var client *elastic.Client

type Elastic struct {
	indexName     string
	elasticClient *elastic.Client
}

func (this *Elastic) Init() {
	if this.elasticClient == nil {
		var err error
		this.elasticClient, err = elastic.NewClient(
			elastic.SetURL(ELASTIC_URL),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Fatal(err)
		}
		this.indexName = ELASTIC_INDEX_NAME
	}
}

func (this *Elastic) Dispose() {
	//return
}

func (this *Elastic) CreateIndex() {
	shards := 1
	if config.ELASTIC_MODE == "cluster" {
		shards = 2
	}

	log.Printf("Create index with shards %d", shards)

	mapping := fmt.Sprintf(`{
		"settings": {
			"number_of_shards": %d,
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
	}`, shards)

	// Create an index with the defined settings and mappings
	createIndex, err := this.elasticClient.CreateIndex(this.indexName).BodyString(mapping).Do(context.Background())
	if err != nil {
		log.Println("Failed to create index: %s", err)
	} else if !createIndex.Acknowledged {
		log.Println("Create index not acknowledged")
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

func (this *Elastic) GetAllDocuments() []driver.Car {
	var response []driver.Car
	response = make([]driver.Car, 0)

	// Initialize scrolling over documents
	scroll := this.elasticClient.Scroll(this.indexName).Size(10000) // Adjust size as needed
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
			var doc driver.Car
			if err := json.Unmarshal(hit.Source, &doc); err != nil {
				log.Fatalf("Error deserializing document: %s", err)
			}
			doc.UUID = hit.Id
			response = append(response, doc)
			// Process your document (doc) here
			log.Printf("Doc ID: %s, Doc: %+v\n", hit.Id, doc)
		}
	}
	return response
}

func (this *Elastic) Search3(regex string) []driver.Car {
	var response []driver.Car
	response = make([]driver.Car, 0)

	fmt.Println(regex)

	query := elastic.NewBoolQuery().Should(
		elastic.NewRegexpQuery("rendszam", ".*"+regex+".*"),
		elastic.NewRegexpQuery("tulajdonos", ".*"+regex+".*"),
		elastic.NewRegexpQuery("adatok", ".*"+regex+".*"),
	)

	searchResult, err := this.elasticClient.Search().
		Index(this.indexName).
		Query(query).
		Pretty(true).
		Size(1000).
		Do(context.Background())

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	fmt.Printf("Found %d documents\n", searchResult.TotalHits())

	for _, hit := range searchResult.Hits.Hits {
		var doc driver.Car
		err := json.Unmarshal(hit.Source, &doc)

		if err != nil {
			log.Fatalf("Error deserializing hit to document: %s", err)
		}

		doc.UUID = hit.Id
		response = append(response, doc)
		fmt.Printf("Document ID: %s, Fields: %+v\n", hit.Id, doc)
	}

	return response
}

func (this *Elastic) AddDocument(doc driver.Car) bool {
	uuid5 := utils.GetUUID(doc.PlateNumber)
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

func (this *Elastic) InsertOne(car driver.Car) *driver.Car {
	uuid5 := utils.GetUUID(car.PlateNumber)
	car.UUID = uuid5

	indexResponse, err := this.elasticClient.Index().
		Index(this.indexName).
		BodyJson(car).
		Id(uuid5).
		Do(context.Background())

	if err != nil {
		log.Println(err)
		return nil
	}

	log.Printf("Indexed document %s to index %s\n", indexResponse.Id, indexResponse.Index)
	return &car
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

func (this *Elastic) CountDocuments() int64 {
	count, err := this.elasticClient.Count().
		Index(this.indexName).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting document count: %s", err)
	}

	fmt.Printf("Document count in '%s': %d\n", this.indexName, count)
	return int64(count)
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

func (this *Elastic) GetByUUID(UUID string) *driver.Car {
	log.Println(UUID)
	response, err := this.elasticClient.Get().
		Index(this.indexName).
		Id(UUID).
		Do(context.Background())

	if err != nil {
		log.Println("Error getting document by ID", err)
		return nil
	}

	if !response.Found {
		log.Printf("Document with ID %s not found in index %s\n", UUID, this.indexName)
		return nil
	}

	var doc driver.Car
	errJson := json.Unmarshal(response.Source, &doc)

	if errJson != nil {
		log.Fatalf("Error deserializing hit to document: %s", err)
		return nil
	}

	return &doc
}

func (this *Elastic) GetIndexStructure() {
	mappings, err := this.elasticClient.GetMapping().Index(this.indexName).Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting mappings: %s", err)
	}

	mappingsJSON, _ := json.MarshalIndent(mappings, "", "  ")
	fmt.Printf("Mappings for index %s: %s\n", this.indexName, string(mappingsJSON))

	settings, err := this.elasticClient.IndexGetSettings(this.indexName).Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting settings: %s", err)
	}

	settingsJSON, _ := json.MarshalIndent(settings, "", "  ")
	fmt.Printf("Settings for index %s: %s\n", this.indexName, string(settingsJSON))
}
