package elasticsearch

import (
	"context"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/go-plugin-utils/utils"
	elastic "gopkg.in/olivere/elastic.v5"
)

// PluginResults a malice plugin results object
type PluginResults struct {
	ID       string `json:"id"`
	Name     string
	Category string
	Data     map[string]interface{}
}

// ElasticAddr ElasticSearch address to user for connections
var ElasticAddr string

// InitElasticSearch initalizes ElasticSearch for use with malice
func InitElasticSearch(elasticHost string) error {

	if ElasticAddr == "" {
		if elasticHost == "" {
			ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elastic"))
		} else {
			ElasticAddr = fmt.Sprintf("http://%s:9200", elasticHost)
		}
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
	if err != nil {
		return err
	}

	exists, err := client.IndexExists("malice").Do(context.Background())
	if err != nil {
		return err
	}

	if !exists {
		// Index does not exist yet.
		log.Debug("Mapping: ", mapping)
		createIndex, err := client.CreateIndex("malice").BodyString(mapping).Do(context.Background())
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			log.Error("Couldn't create Index.")
		} else {
			log.Debug("Created Index: ", "malice")
		}
	} else {
		log.Debug("Index malice already exists.")
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

// WritePluginResultsToDatabase upserts plugin results into Database
func WritePluginResultsToDatabase(results PluginResults) error {
	// log.Info(results)
	// scanID := utils.Getopt("MALICE_SCANID", "")
	if ElasticAddr == "" {
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elastic"))
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
	if err != nil {
		return err
	}

	getSample, err := client.Get().
		Index("malice").
		Type("samples").
		Id(results.ID).
		Do(context.Background())
	if err != nil {
		log.Debug(err)
	}
	// utils.Assert(err)

	if getSample != nil && getSample.Found {
		log.Debugf("Got document %s in version %d from index %s, type %s\n", getSample.Id, getSample.Version, getSample.Index, getSample.Type)
		updateScan := map[string]interface{}{
			"scan_date": time.Now().Format(time.RFC3339Nano),
			"plugins": map[string]interface{}{
				results.Category: map[string]interface{}{
					results.Name: results.Data,
				},
			},
		}
		update, err := client.Update().Index("malice").Type("samples").Id(getSample.Id).
			Doc(updateScan).
			Do(context.Background())
		if err != nil {
			return err
		}

		log.Debugf("New version of sample %q is now %d\n", update.Id, update.Version)
		// return *update

	} else {

		scan := map[string]interface{}{
			// "id":      sample.SHA256,
			// "file":      sample,
			"plugins": map[string]interface{}{
				results.Category: map[string]interface{}{
					results.Name: results.Data,
				},
			},
			"scan_date": time.Now().Format(time.RFC3339Nano),
		}

		newScan, err := client.Index().
			Index("malice").
			Type("samples").
			OpType("create").
			// Id("1").
			BodyJson(scan).
			Do(context.Background())
		if err != nil {
			return err
		}

		log.Debugf("Indexed sample %s to index %s, type %s\n", newScan.Id, newScan.Index, newScan.Type)
		// return *newScan
	}
	return err
}
