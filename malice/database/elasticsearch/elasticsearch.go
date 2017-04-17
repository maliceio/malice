package elasticsearch

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	mallog "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/database"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/maldirs"
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

func getElasticSearchAddr(addr string) string {
	if _, inDocker := os.LookupEnv("MALICE_IN_DOCKER"); inDocker {
		if addr != "" {
			return fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", addr))
		}
		return fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elasticsearch"))
	}
	return fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "localhost"))
}

// Start creates an Elasticsearch container from the image blacktop/elasticsearch
func Start(docker *client.Docker, logs bool) (types.ContainerJSONBase, error) {

	name := config.Conf.DB.Name
	image := config.Conf.DB.Image
	binds := []string{"malice:/usr/share/elasticsearch/data"}
	portBindings := nat.PortMap{
		"9200/tcp": {{HostIP: "0.0.0.0", HostPort: "9200"}},
	}

	if docker.Ping() {
		cont, err := container.Start(docker, nil, name, image, logs, binds, portBindings, nil, nil)
		if err != nil {
			return types.ContainerJSONBase{}, err
		}
		// Inspect newly created container to get IP assigned to it
		dbInfo, err := container.Inspect(docker, cont.ID)
		if err != nil {
			mallog.Error(err)
		}
		elasticAddress := getElasticSearchAddr(dbInfo.NetworkSettings.IPAddress)

		mallog.WithFields(mallog.Fields{
			// "id":   cont.ID,
			"ip":   docker.GetIP(),
			"port": config.Conf.DB.Ports,
			"name": cont.Name,
			"env":  config.Conf.Environment.Run,
		}).Info("Elasticsearch Container Started")

		// Give ELK a few seconds to start
		mallog.WithFields(mallog.Fields{
			"server":  elasticAddress,
			"timeout": config.Conf.DB.Timeout,
		}).Info("Waiting for Elasticsearch to come online.")

		ctx := context.Background()
		err = WaitForConnection(ctx, "", config.Conf.DB.Timeout)
		if err != nil {
			logOpts := types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Follow:     false,
			}

			logs, _ := docker.Client.ContainerLogs(context.Background(), cont.ID, logOpts)
			defer logs.Close()
			// Convert logs to a string
			buf := new(bytes.Buffer)
			buf.ReadFrom(logs)
			logStr := buf.String()
			mallog.Debug(logStr)
			// Check if elasticsearch could not start due to lack of RAM
			if strings.Contains(logStr, "There is insufficient memory for the Java Runtime Environment to continue") {
				info, err := docker.Client.Info(context.Background())
				if err != nil {
					mallog.Error(err)
				}
				log.Fatal("You do not have enough RAM to run elasticsearch. Elasticsearch needs at least 2GB and you have: ", units.BytesSize(float64(info.MemTotal)))
			}
			if err != nil {
				mallog.Error(err)
			}
		}

		return cont, err
	}
	return types.ContainerJSONBase{}, errors.New("Cannot connect to the Docker daemon. Is the docker daemon running on this host?")
}

// InitElasticSearch initalizes ElasticSearch for use with malice
func InitElasticSearch(addr string) error {

	// Test connection to ElasticSearch
	_, err := TestConnection(addr)
	er.CheckError(err)

	file, err := os.OpenFile(path.Join(maldirs.GetLogsDir(), "elastic.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(ElasticAddr),
		elastic.SetErrorLog(log.New(file, "ERROR ", log.LstdFlags)),
	)
	utils.Assert(err)

	exists, err := client.IndexExists("malice").Do(context.Background())
	utils.Assert(err)

	if !exists {
		// Index does not exist yet.
		createIndex, err := client.CreateIndex("malice").BodyString(mapping).Do(context.Background())
		utils.Assert(err)
		if !createIndex.Acknowledged {
			// Not acknowledged
			mallog.Error("Couldn't create Index.")
		} else {
			mallog.Debug("Created Index: ", "malice")
		}
	} else {
		mallog.Debug("Index malice already exists.")
	}

	return err
}

// TestConnection tests the ElasticSearch connection
func TestConnection(addr string) (bool, error) {

	var err error

	if ElasticAddr == "" {
		ElasticAddr = getElasticSearchAddr(addr)
	}

	// connect to ElasticSearch where --link elastic was using via malice in Docker
	file, err := os.OpenFile(path.Join(maldirs.GetLogsDir(), "elastic.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(ElasticAddr),
		elastic.SetErrorLog(log.New(file, "ERROR ", log.LstdFlags)),
	)
	if err != nil {
		return false, err
	}

	// Ping the Elasticsearch server to get e.g. the version number
	mallog.Debugf("Attempting to PING to: %s", ElasticAddr)
	info, code, err := client.Ping(ElasticAddr).Do(context.Background())
	if err != nil {
		return false, err
	}

	mallog.WithFields(mallog.Fields{
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

// WaitForConnection waits for connection to Elasticsearch to be ready
func WaitForConnection(ctx context.Context, addr string, timeout int) error {

	var ready bool
	var connErr error
	secondsWaited := 0

	connCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	mallog.Debug("===> trying to connect to elasticsearch")
	for {
		// Try to connect to Elasticsearch
		select {
		case <-connCtx.Done():
			mallog.WithFields(mallog.Fields{"timeout": timeout}).Error("connecting to elasticsearch timed out")
			return connErr
		default:
			ready, connErr = TestConnection(addr)
			if ready {
				mallog.Infof("Elasticsearch came online after %d seconds", secondsWaited)
				return connErr
			}
			secondsWaited++
			time.Sleep(1 * time.Second)
		}
	}
}

// WriteFileToDatabase inserts sample into Database
func WriteFileToDatabase(sample map[string]interface{}) elastic.IndexResponse {

	// Test connection to ElasticSearch
	_, err := TestConnection("")
	er.CheckError(err)

	file, err := os.OpenFile(path.Join(maldirs.GetLogsDir(), "elastic.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(ElasticAddr),
		elastic.SetErrorLog(log.New(file, "ERROR ", log.LstdFlags)),
	)
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
		OpType("index").
		// Id("1").
		BodyJson(scan).
		Do(context.Background())
	utils.Assert(err)

	mallog.WithFields(mallog.Fields{
		"id":    newScan.Id,
		"index": newScan.Index,
		"type":  newScan.Type,
	}).Debug("Indexed sample.")

	return *newScan
}

// WriteHashToDatabase inserts sample into Database
func WriteHashToDatabase(hash string) elastic.IndexResponse {

	hashType, err := utils.GetHashType(hash)
	utils.Assert(err)

	file, err := os.OpenFile(path.Join(maldirs.GetLogsDir(), "elastic.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(ElasticAddr),
		elastic.SetErrorLog(log.New(file, "ERROR ", log.LstdFlags)),
	)
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
		Do(context.Background())
	utils.Assert(err)

	mallog.WithFields(mallog.Fields{
		"id":    newScan.Id,
		"index": newScan.Index,
		"type":  newScan.Type,
	}).Debug("Indexed sample.")

	return *newScan
}

// WritePluginResultsToDatabase upserts plugin results into Database
func WritePluginResultsToDatabase(results PluginResults) {

	// scanID := utils.Getopt("MALICE_SCANID", "")
	if ElasticAddr == "" {
		ElasticAddr = fmt.Sprintf("http://%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elastic"))
	}

	file, err := os.OpenFile(path.Join(maldirs.GetLogsDir(), "elastic.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(ElasticAddr),
		elastic.SetErrorLog(log.New(file, "ERROR ", log.LstdFlags)),
	)
	utils.Assert(err)

	getSample, err := client.Get().
		Index("malice").
		Type("samples").
		Id(results.ID).
		Do(context.Background())

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
			Do(context.Background())
		utils.Assert(err)

		mallog.WithFields(mallog.Fields{
			"id":      update.Id,
			"version": update.Version,
		}).Debug("New version of sample.")

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
		utils.Assert(err)

		mallog.WithFields(mallog.Fields{
			"id":    newScan.Id,
			"index": newScan.Index,
			"type":  newScan.Type,
		}).Debug("Indexed sample.")
		// return *newScan
	}
}
