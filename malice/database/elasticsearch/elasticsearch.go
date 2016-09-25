package elasticsearch

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/maliceio/malice/malice/database"
	"github.com/maliceio/malice/malice/persist"
	util "github.com/maliceio/malice/utils"
	elastic "gopkg.in/olivere/elastic.v3"
)

// ElasticAddr ElasticSearch address to user for connections
var ElasticAddr string

// InitElasticSearch initalizes ElasticSearch for use with malice
func InitElasticSearch() error {
	client, err := elastic.NewSimpleClient()
	utils.Assert(err)

	exists, err := client.IndexExists("malice").Do()
	utils.Assert(err)

	if !exists {
		// Index does not exist yet.
		createIndex, err := client.CreateIndex("malice").BodyString(mapping).Do()
		utils.Assert(err)
		if !createIndex.Acknowledged {
			// Not acknowledged
			log.Error("Couldn't create Index.")
		} else {
			log.Info("Created Index: ", "malice")
		}
	} else {
		log.Info("Index malice already exists.")
	}

	return err
}

// TestConnection tests the ElasticSearch connection
func TestConnection(addr string) error {

	if ElasticAddr == "" {
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elastic"))
	}

	// connect to ElasticSearch where --link elastic was using via malice in Docker
	log.Debugf("Attempting to connect to: %s", ElasticAddr)
	_, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))

	// Ping the Elasticsearch server to get e.g. the version number
	// info, code, err := client.Ping(ElasticAddr).Do()
	// utils.Assert(err)
	// fmt.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

	if err != nil {
		// connect to ElasticSearch via malice in Docker
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", addr))
		log.Debugf("Attempting to connect to: %s", ElasticAddr)
		_, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))

		if err != nil {
			// connect to ElasticSearch using Docker for Mac
			ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "localhost"))
			log.Debugf("Attempting to connect to: %s", ElasticAddr)
			_, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
			return err
		}
		return err
	}
	return err
}

// WriteFileToDatabase inserts sample into Database
func WriteFileToDatabase(sample persist.File) elastic.IndexResponse {
	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
	utils.Assert(err)

	// getSample, err := client.Get().
	// 	Index("malice").
	// 	Type("samples").
	// 	Id("1").
	// 	Do()

	// fmt.Println(getSample)
	// fmt.Println(err)
	// if err != nil {

	// }

	// if getSample.Found {
	// 	fmt.Printf("Got document %s in version %d from index %s, type %s\n", getSample.Id, getSample.Version, getSample.Index, getSample.Type)
	// } else {

	scan := map[string]interface{}{
		// "id":      sample.SHA256,
		"file":      sample,
		"plugins":   database.GetPluginsByCategory(),
		"scan_date": time.Now().Format(time.RFC3339Nano),
	}

	newScan, err := client.Index().
		Index("malice").
		Type("samples").
		OpType("create").
		// Id("1").
		BodyJson(scan).
		Do()
	utils.Assert(err)
	log.Debugf("Indexed sample %s to index %s, type %s\n", newScan.Id, newScan.Index, newScan.Type)

	update, err := client.Update().Index("malice").Type("samples").Id(newScan.Id).
		Doc(map[string]interface{}{
			"plugins": map[string]interface{}{
				"intel": map[string]interface{}{
					"nsrl": "UPDATED",
				},
			},
		}).
		Do()
	utils.Assert(err)
	log.Debugf("New version of sample %q is now %d\n", update.Id, update.Version)

	// }

	return *newScan
}

// WriteHashToDatabase inserts sample into Database
func WriteHashToDatabase(hash string) elastic.IndexResponse {

	hashType, err := util.GetHashType(hash)
	utils.Assert(err)

	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
	utils.Assert(err)

	scan := map[string]interface{}{
		// "id":      sample.SHA256,
		"file": map[string]interface{}{
			hashType: hash,
		},
		"plugins":   database.GetPluginsByCategory(),
		"scan_date": time.Now().Format(time.RFC3339Nano),
	}

	newScan, err := client.Index().
		Index("malice").
		Type("samples").
		OpType("create").
		// Id("1").
		BodyJson(scan).
		Do()
	utils.Assert(err)
	log.Debugf("Indexed sample %s to index %s, type %s\n", newScan.Id, newScan.Index, newScan.Type)

	return *newScan
}
