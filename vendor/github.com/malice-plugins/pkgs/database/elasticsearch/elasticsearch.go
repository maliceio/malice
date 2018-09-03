package elasticsearch

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/malice-plugins/pkgs/database"
	"github.com/malice-plugins/pkgs/utils"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

// Database is the elasticsearch malice database object
type Database struct {
	Host     string                 `json:"host,omitempty"`
	Port     string                 `json:"port,omitempty"`
	URL      string                 `json:"url,omitempty"`
	Username string                 `json:"username,omitempty"`
	Password string                 `json:"password,omitempty"`
	Index    string                 `json:"index,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Plugins  map[string]interface{} `json:"plugins,omitempty"`
}

var (
	defaultIndex string
	defaultType  string
	defaultHost  string
	defaultPort  string
	defaultURL   string
)

func init() {
	defaultIndex = utils.Getopt("MALICE_ELASTICSEARCH_INDEX", "malice")
	defaultType = utils.Getopt("MALICE_ELASTICSEARCH_TYPE", "samples")
	defaultHost = utils.Getopt("MALICE_ELASTICSEARCH_HOST", "localhost")
	defaultPort = utils.Getopt("MALICE_ELASTICSEARCH_PORT", "9200")
}

// getURL with the following order of precedence
// - user input (cli)
// - user ENV
// - sane defaults
func (db *Database) getURL() {

	// If not set use defaults
	if len(strings.TrimSpace(db.Index)) == 0 {
		db.Index = defaultIndex
	}
	if len(strings.TrimSpace(db.Type)) == 0 {
		db.Type = defaultType
	}
	if len(strings.TrimSpace(db.Host)) == 0 {
		db.Host = defaultHost
	}
	if len(strings.TrimSpace(db.Port)) == 0 {
		db.Port = defaultPort
	}

	// If user set URL param use it
	if len(strings.TrimSpace(db.URL)) == 0 {
		// If running in docker use `elasticsearch`
		if _, exists := os.LookupEnv("MALICE_IN_DOCKER"); exists {
			db.URL = utils.Getopt("MALICE_ELASTICSEARCH_URL", fmt.Sprintf("%s:%s", "elasticsearch", db.Port))
			log.WithField("elasticsearch_url", db.URL).Debug("running malice in docker")
			return
		}

		db.URL = utils.Getopt("MALICE_ELASTICSEARCH_URL", fmt.Sprintf("%s:%s", db.Host, db.Port))
	}
}

// Init initalizes ElasticSearch for use with malice
func (db *Database) Init() error {

	// Create URL from host/port
	db.getURL()

	// Test connection to ElasticSearch
	err := db.TestConnection()
	if err != nil {
		return errors.Wrap(err, "failed to connect to database")
	}

	client, err := elastic.NewSimpleClient(
		elastic.SetURL(db.URL),
		elastic.SetBasicAuth(
			utils.Getopts(db.Username, "MALICE_ELASTICSEARCH_USERNAME", ""),
			utils.Getopts(db.Password, "MALICE_ELASTICSEARCH_PASSWORD", ""),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create elasticsearch simple client")
	}

	exists, err := client.IndexExists(db.Index).Do(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to check if index exists")
	}

	if !exists {
		// Index does not exist yet.
		createIndex, err := client.CreateIndex(db.Index).BodyString(mapping).Do(context.Background())
		if err != nil {
			return errors.Wrapf(err, "failed to create index: %s", db.Index)
		}

		if !createIndex.Acknowledged {
			log.Error("index creation not acknowledged")
		} else {
			log.Debugf("created index %s", db.Index)
		}
	} else {
		log.Debugf("index %s already exists", db.Index)
	}

	return nil
}

// TestConnection tests the ElasticSearch connection
func (db *Database) TestConnection() error {

	// Create URL from host/port
	db.getURL()

	// connect to ElasticSearch where --link elasticsearch was using via malice in Docker
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(db.URL),
		elastic.SetBasicAuth(
			utils.Getopts(db.Username, "MALICE_ELASTICSEARCH_USERNAME", ""),
			utils.Getopts(db.Password, "MALICE_ELASTICSEARCH_PASSWORD", ""),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create elasticsearch simple client")
	}

	// Ping the Elasticsearch server to get e.g. the version number
	log.Debugf("attempting to PING to: %s", db.URL)
	info, code, err := client.Ping(db.URL).Do(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to ping elasticsearch")
	}

	log.WithFields(log.Fields{
		"code":    code,
		"cluster": info.ClusterName,
		"version": info.Version.Number,
		"url":     db.URL,
	}).Debug("elasticSearch connection successful")

	return nil
}

// WaitForConnection waits for connection to Elasticsearch to be ready
func (db *Database) WaitForConnection(ctx context.Context, timeout int) error {

	var err error

	secondsWaited := 0

	connCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	log.Debug("===> trying to connect to elasticsearch")
	for {
		// Try to connect to Elasticsearch
		select {
		case <-connCtx.Done():
			return errors.Wrapf(err, "connecting to elasticsearch timed out after %d seconds", secondsWaited)
		default:
			err = db.TestConnection()
			if err == nil {
				log.Debugf("elasticsearch came online after %d seconds", secondsWaited)
				return nil
			}
			// not ready yet
			secondsWaited++
			log.Debug(" * could not connect to elasticsearch (sleeping for 1 second)")
			time.Sleep(1 * time.Second)
		}
	}
}

// StoreFileInfo inserts initial sample info into database creating a placeholder for it
func (db *Database) StoreFileInfo(sample map[string]interface{}) (elastic.IndexResponse, error) {

	if len(db.Plugins) == 0 {
		return elastic.IndexResponse{}, errors.New("Database.Plugins is empty (you must set this field to use this function)")
	}

	// Test connection to ElasticSearch
	err := db.TestConnection()
	if err != nil {
		return elastic.IndexResponse{}, errors.Wrap(err, "failed to connect to database")
	}

	client, err := elastic.NewSimpleClient(
		elastic.SetURL(db.URL),
		elastic.SetBasicAuth(
			utils.Getopts(db.Username, "MALICE_ELASTICSEARCH_USERNAME", ""),
			utils.Getopts(db.Password, "MALICE_ELASTICSEARCH_PASSWORD", ""),
		),
	)
	if err != nil {
		return elastic.IndexResponse{}, errors.Wrap(err, "failed to create elasticsearch simple client")
	}

	// NOTE: I am not setting ID because I want to be able to re-scan files with updated signatures in the future
	fInfo := map[string]interface{}{
		// "id":      sample.SHA256,
		"file":      sample,
		"plugins":   db.Plugins,
		"scan_date": time.Now().Format(time.RFC3339Nano),
	}

	newScan, err := client.Index().
		Index(db.Index).
		Type(db.Type).
		OpType("index").
		// Id("1").
		BodyJson(fInfo).
		Do(context.Background())
	if err != nil {
		return elastic.IndexResponse{}, errors.Wrap(err, "failed to index file info")
	}

	log.WithFields(log.Fields{
		"id":    newScan.Id,
		"index": newScan.Index,
		"type":  newScan.Type,
	}).Debug("indexed sample")

	return *newScan, nil
}

// StoreHash stores a hash into the database that has been queried via intel-plugins
func (db *Database) StoreHash(hash string) (elastic.IndexResponse, error) {

	if len(db.Plugins) == 0 {
		return elastic.IndexResponse{}, errors.New("Database.Plugins is empty (you must set this field to use this function)")
	}

	hashType, err := utils.GetHashType(hash)
	if err != nil {
		return elastic.IndexResponse{}, errors.Wrapf(err, "unable to detect hash type: %s", hash)
	}

	// Test connection to ElasticSearch
	err = db.TestConnection()
	if err != nil {
		return elastic.IndexResponse{}, errors.Wrap(err, "failed to connect to database")
	}

	client, err := elastic.NewSimpleClient(
		elastic.SetURL(db.URL),
		elastic.SetBasicAuth(
			utils.Getopts(db.Username, "MALICE_ELASTICSEARCH_USERNAME", ""),
			utils.Getopts(db.Password, "MALICE_ELASTICSEARCH_PASSWORD", ""),
		),
	)
	if err != nil {
		return elastic.IndexResponse{}, errors.Wrap(err, "failed to create elasticsearch simple client")
	}

	scan := map[string]interface{}{
		// "id":      sample.SHA256,
		"file": map[string]interface{}{
			hashType: hash,
		},
		"plugins":   db.Plugins,
		"scan_date": time.Now().Format(time.RFC3339Nano),
	}

	newScan, err := client.Index().
		Index(db.Index).
		Type(db.Type).
		OpType("create").
		// Id("1").
		BodyJson(scan).
		Do(context.Background())
	if err != nil {
		return elastic.IndexResponse{}, errors.Wrapf(err, "unable to index hash: %s", hash)
	}

	log.WithFields(log.Fields{
		"id":    newScan.Id,
		"index": newScan.Index,
		"type":  newScan.Type,
	}).Debug("indexed sample")

	return *newScan, nil
}

// StorePluginResults stores a plugin's results in the database by updating
// the placeholder created by the call to StoreFileInfo
func (db *Database) StorePluginResults(results database.PluginResults) error {

	// Test connection to ElasticSearch
	err := db.TestConnection()
	if err != nil {
		return errors.Wrap(err, "failed to connect to database")
	}

	client, err := elastic.NewSimpleClient(
		elastic.SetURL(db.URL),
		elastic.SetBasicAuth(
			utils.Getopts(db.Username, "MALICE_ELASTICSEARCH_USERNAME", ""),
			utils.Getopts(db.Password, "MALICE_ELASTICSEARCH_PASSWORD", ""),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create elasticsearch simple client")
	}

	// get sample db record
	getSample, err := client.Get().
		Index(db.Index).
		Type(db.Type).
		Id(results.ID).
		Do(context.Background())
	// ignore 404 not found error
	if err != nil && !elastic.IsNotFound(err) {
		return errors.Wrapf(err, "failed to get sample with id: %s", results.ID)
	}

	if getSample != nil && getSample.Found {
		log.Debugf("got document %s in version %d from index %s, type %s\n", getSample.Id, getSample.Version, getSample.Index, getSample.Type)
		updateScan := map[string]interface{}{
			"scan_date": time.Now().Format(time.RFC3339Nano),
			"plugins": map[string]interface{}{
				results.Category: map[string]interface{}{
					results.Name: results.Data,
				},
			},
		}
		update, err := client.Update().Index(db.Index).Type(db.Type).Id(getSample.Id).
			Doc(updateScan).
			Do(context.Background())
		if err != nil {
			return errors.Wrapf(err, "failed to update sample with id: %s", results.ID)
		}

		log.Debugf("updated version of sample %q is now %d\n", update.Id, update.Version)

	} else {
		// ID not found so create new document with `index` command
		scan := map[string]interface{}{
			"plugins": map[string]interface{}{
				results.Category: map[string]interface{}{
					results.Name: results.Data,
				},
			},
			"scan_date": time.Now().Format(time.RFC3339Nano),
		}

		newScan, err := client.Index().
			Index(db.Index).
			Type(db.Type).
			OpType("index").
			// Id("1").
			BodyJson(scan).
			Do(context.Background())
		if err != nil {
			return errors.Wrapf(err, "failed to create new sample plugin doc with id: %s", results.ID)
		}

		log.WithFields(log.Fields{
			"id":    newScan.Id,
			"index": newScan.Index,
			"type":  newScan.Type,
		}).Debug("indexed sample")
	}

	return nil
}
