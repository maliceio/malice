package elasticsearch

import (
	"errors"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/go-connections/nat"

	"github.com/maliceio/malice/malice/database"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	"github.com/maliceio/malice/utils"
	elastic "gopkg.in/olivere/elastic.v3"
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

// StartELK creates an ELK container from the image blacktop/elk
func StartELK(docker *client.Docker, logs bool) (types.ContainerJSONBase, error) {

	name := "elk"
	image := "blacktop/elk"
	binds := []string{"malice:/usr/share/elasticsearch/data"}
	portBindings := nat.PortMap{
		"80/tcp":   {{HostIP: "0.0.0.0", HostPort: "80"}},
		"9200/tcp": {{HostIP: "0.0.0.0", HostPort: "9200"}},
	}

	if docker.Ping() {
		cont, err := container.Start(docker, nil, name, image, logs, binds, portBindings, nil, nil)

		// Give ELK a few seconds to start
		time.Sleep(10 * time.Second)
		log.Info("sleeping for 5 seconds to let ELK start")
		return cont, err
	}
	return types.ContainerJSONBase{}, errors.New("Cannot connect to the Docker daemon. Is the docker daemon running on this host?")
}

// InitElasticSearch initalizes ElasticSearch for use with malice
func InitElasticSearch() error {

	if ElasticAddr == "" {
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elk"))
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
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
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "localhost"))
	}

	// connect to ElasticSearch where --link elastic was using via malice in Docker
	log.Debugf("Attempting to connect to: %s", ElasticAddr)
	_, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))

	// Ping the Elasticsearch server to get e.g. the version number
	// info, code, err := client.Ping(ElasticAddr).Do()
	// utils.Assert(err)
	// fmt.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

	// if err != nil {
	// 	// connect to ElasticSearch via malice in Docker
	// 	ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elk"))
	// 	log.Debugf("Attempting to connect to: %s", ElasticAddr)
	// 	_, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))

	if err != nil {
		// connect to ElasticSearch using Docker for Mac
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", addr))
		log.Debugf("Attempting to connect to: %s", ElasticAddr)
		_, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
		return err
	}
	// return err
	// }
	return err
}

// WriteFileToDatabase inserts sample into Database
func WriteFileToDatabase(sample map[string]interface{}) elastic.IndexResponse {

	if ElasticAddr == "" {
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elastic"))
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
	utils.Assert(err)

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

	// update, err := client.Update().Index("malice").Type("samples").Id(newScan.Id).
	// 	Doc(map[string]interface{}{
	// 		"plugins": map[string]interface{}{
	// 			"intel": map[string]interface{}{
	// 				"nsrl": "UPDATED",
	// 			},
	// 		},
	// 	}).
	// 	Do()
	// utils.Assert(err)
	// log.Debugf("New version of sample %q is now %d\n", update.Id, update.Version)

	return *newScan
}

// WriteHashToDatabase inserts sample into Database
func WriteHashToDatabase(hash string) elastic.IndexResponse {

	hashType, err := utils.GetHashType(hash)
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

// WritePluginResultsToDatabase upserts plugin results into Database
func WritePluginResultsToDatabase(results PluginResults) {

	// scanID := utils.Getopt("MALICE_SCANID", "")
	if ElasticAddr == "" {
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elastic"))
	}

	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
	utils.Assert(err)

	getSample, err := client.Get().
		Index("malice").
		Type("samples").
		Id(results.ID).
		Do()

	if err != nil {

	}

	if getSample != nil && getSample.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", getSample.Id, getSample.Version, getSample.Index, getSample.Type)
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
			Do()
		utils.Assert(err)

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
			Do()
		utils.Assert(err)

		log.Debugf("Indexed sample %s to index %s, type %s\n", newScan.Id, newScan.Index, newScan.Type)
		// return *newScan
	}
}
