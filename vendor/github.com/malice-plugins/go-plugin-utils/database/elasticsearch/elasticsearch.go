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
func InitElasticSearch(addr string) error {

	// Test connection to ElasticSearch
	_, err := TestConnection(addr)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "InitElasticSearch.TestConnection",
		}).Debug(err)
		return err
	}

	client, err := elastic.NewSimpleClient(
		elastic.SetURL(ElasticAddr),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "InitElasticSearch.NewSimpleClient",
		}).Debug(err)
		return err
	}

	exists, err := client.IndexExists("malice").Do(context.Background())
	if err != nil {
		log.WithFields(log.Fields{
			"func": "InitElasticSearch.IndexExists",
		}).Debug(err)
	}

	if !exists {
		// Index does not exist yet.
		createIndex, err := client.CreateIndex("malice").BodyString(mapping).Do(context.Background())
		if err != nil {
			log.WithFields(log.Fields{
				"func": "InitElasticSearch.CreateIndex",
			}).Debug(err)
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
func TestConnection(addr string) (bool, error) {

	var err error

	if ElasticAddr == "" {
		if addr == "" {
			ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elasticsearch"))
		} else {
			ElasticAddr = addr
		}
	}

	// connect to ElasticSearch where --link elastic was using via malice in Docker
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(ElasticAddr),
	)
	if err != nil {
		return false, err
	}

	// Ping the Elasticsearch server to get e.g. the version number
	log.Debugf("Attempting to PING to: %s", ElasticAddr)
	info, code, err := client.Ping(ElasticAddr).Do(context.Background())
	if err != nil {
		return false, err
	}

	log.WithFields(log.Fields{
		"code":    code,
		"cluster": info.ClusterName,
		"version": info.Version.Number,
		"address": ElasticAddr,
	}).Debug("ElasticSearch connection successful.")

	if code == 200 {
		return true, err
	}
	return false, err
}

// WritePluginResultsToDatabase upserts plugin results into Database
func WritePluginResultsToDatabase(results PluginResults) error {
	// log.Info(results)
	// scanID := utils.Getopt("MALICE_SCANID", "")
	if ElasticAddr == "" {
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elasticsearch"))
		log.Debug("Using elasticsearch address: ", ElasticAddr)
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
	if err != nil {
		log.WithFields(log.Fields{
			"func": "WritePluginResultsToDatabase.NewSimpleClient",
		}).Debug(err)
		return err
	}

	getSample, err := client.Get().
		Index("malice").
		Type("samples").
		Id(results.ID).
		Do(context.Background())
	if err != nil {
		log.WithFields(log.Fields{
			"id":    results.ID,
			"index": "malice",
			"type":  "samples",
		}).Debug("Trying to find document and got error: ", err)
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
			log.WithFields(log.Fields{
				"func": "WritePluginResultsToDatabase.Update",
			}).Debug(err)
			return err
		}

		log.Debugf("New version of sample %q is now %d\n", update.Id, update.Version)
		// return *update
	} else {
		// ID not found so create new document with `index` command
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
			OpType("index").
			// Id("1").
			BodyJson(scan).
			Do(context.Background())
		if err != nil {
			log.WithFields(log.Fields{
				"func": "WritePluginResultsToDatabase.Index",
			}).Debug(err)
			return err
		}

		log.WithFields(log.Fields{
			"id":    newScan.Id,
			"index": newScan.Index,
			"type":  newScan.Type,
		}).Debug("Indexed sample.")
	}

	return err
}
