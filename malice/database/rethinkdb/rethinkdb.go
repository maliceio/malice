package rethinkdb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/maliceio/go-plugin-utils/waitforit"
	"github.com/maliceio/malice/config"
	"github.com/maliceio/malice/malice/database"
	"github.com/maliceio/malice/malice/docker/client"
	"github.com/maliceio/malice/malice/docker/client/container"
	er "github.com/maliceio/malice/malice/errors"
	"github.com/maliceio/malice/malice/persist"
	r "gopkg.in/dancannon/gorethink.v2"
)

// RethinkAddr RethinkDB address to user for connections
var RethinkAddr string

func init() {
	r.Log.Out = ioutil.Discard
}

// StartRethinkDB creates an RethinkDB container from the image rethinkdb
func StartRethinkDB(docker *client.Docker, logs bool) (types.ContainerJSONBase, error) {

	name := config.Conf.DB.Name
	image := "rethinkdb"
	binds := []string{"malice:/data"}
	portBindings := nat.PortMap{
		"8080/tcp":  {{HostIP: "0.0.0.0", HostPort: "8081"}},
		"28015/tcp": {{HostIP: "0.0.0.0", HostPort: "28015"}},
	}

	if docker.Ping() {
		cont, err := container.Start(docker, nil, name, image, logs, binds, portBindings, nil, nil)
		// er.CheckError(err)
		// if network, exists, _ := docker.NetworkExists("malice"); exists {
		// 	err := docker.ConnectNetwork(network, cont)
		// 	er.CheckError(err)
		// }

		// Give rethinkDB a few seconds to start
		log.WithFields(log.Fields{
			"server":  config.Conf.DB.Server,
			"port":    config.Conf.DB.Ports[0],
			"timeout": config.Conf.DB.Timeout,
		}).Debug("Waiting for RethinkDB to come online.")
		er.CheckError(waitforit.WaitForIt(
			"", // fullConn string,
			config.Conf.DB.Server,   // host string,
			config.Conf.DB.Ports[0], // port int,
			config.Conf.DB.Timeout,  // timeout int,
		))
		log.Debug("RethinkDB is now online.")

		return cont, err
	}
	return types.ContainerJSONBase{}, errors.New("Cannot connect to the Docker daemon. Is the docker daemon running on this host?")
}

// TestConnection tests the rethinkDB connection
func TestConnection(addr string) error {

	if RethinkAddr == "" {
		RethinkAddr = fmt.Sprintf("%s:28015", utils.Getopt("MALICE_RETHINKDB", "rethink"))
	}

	// connect to RethinkDB where --link rethink was using via malice in Docker
	log.Debugf("Attempting to connect to: %s", RethinkAddr)
	session, err := r.Connect(r.ConnectOpts{
		Address: RethinkAddr,
	})
	if err != nil {
		// connect to RethinkDB via malice in Docker
		RethinkAddr = fmt.Sprintf("%s:28015", utils.Getopt("MALICE_RETHINKDB", addr))
		log.Debugf("Attempting to connect to: %s", RethinkAddr)
		session, err := r.Connect(r.ConnectOpts{
			Address: RethinkAddr,
			Timeout: 2 * time.Second,
		})
		if err != nil {
			// connect to RethinkDB using Docker for Mac
			RethinkAddr = fmt.Sprintf("%s:28015", utils.Getopt("MALICE_RETHINKDB", "localhost"))
			log.Debugf("Attempting to connect to: %s", RethinkAddr)
			session, err := r.Connect(r.ConnectOpts{
				Address: RethinkAddr,
			})
			defer session.Close()
			return err
		}
		defer session.Close()
		return err
	}
	defer session.Close()
	return err
}

// InitRethinkDB initalizes rethinkDB for use with malice
func InitRethinkDB() error {

	// connect to RethinkDB
	log.Debugf("Attempting to connect to: %s", RethinkAddr)
	session, err := r.Connect(r.ConnectOpts{
		Address: RethinkAddr,
	})
	defer session.Close()
	utils.Assert(err)
	// Delete database test if it exists
	resp, err := r.DBDrop("test").RunWrite(session)
	if err != nil {
		log.Debug(err)
	} else {
		log.Infof("%d DB deleted\n", resp.DBsDropped)
	}
	// Create malice DB
	resp, err = r.DBCreate("malice").RunWrite(session)
	if err != nil {
		log.Debug(err)
	} else {
		log.Infof("%d DB created\n", resp.DBsCreated)
	}
	// Create channel Table
	resp, err = r.DB("malice").TableCreate("samples").RunWrite(session)
	if err != nil {
		log.Debug(err)
	} else {
		log.Infof("%d Table created\n", resp.TablesCreated)
	}
	session.Use("malice")

	return err
}

// WriteFileToDatabase inserts sample into Database
func WriteFileToDatabase(sample persist.File) r.WriteResponse {

	// connect to RethinkDB
	log.Debugf("Attempting to connect to: %s", RethinkAddr)
	session, err := r.Connect(r.ConnectOpts{
		Address:  RethinkAddr,
		Database: "malice",
	})
	defer session.Close()

	if err == nil {
		res, err := r.Table("samples").Filter(map[string]interface{}{
			"file": map[string]interface{}{
				"sha256": sample.SHA256,
			},
		}).Run(session)
		utils.Assert(err)
		defer res.Close()

		var samples []interface{}
		err = res.All(&samples)
		if err != nil {
			log.Errorf("Error scanning database result: %s\n", err)
			return r.WriteResponse{}
		}

		log.Debugf("%d samples\n", len(samples))
		log.Debugln("res: ", res)

		if res.IsNil() {
			// upsert into RethinkDB
			resp, err := r.Table("samples").Insert(
				map[string]interface{}{
					// "id":      sample.SHA256,
					"file":    sample,
					"plugins": database.GetPluginsByCategory(),
				}, r.InsertOpts{Conflict: "replace"}).RunWrite(session)
			utils.Assert(err)

			return resp
		}
		log.Debugf("Sample: %s already exists in the database.", sample.SHA256)
		return r.WriteResponse{}

	}
	log.Debug(err)

	return r.WriteResponse{}
}

// WriteHashToDatabase inserts hash into Database
func WriteHashToDatabase(hash string) r.WriteResponse {

	hashType, err := utils.GetHashType(hash)
	utils.Assert(err)

	// connect to RethinkDB
	log.Debugf("Attempting to connect to: %s", RethinkAddr)
	session, err := r.Connect(r.ConnectOpts{
		Address:  RethinkAddr,
		Database: "malice",
	})
	defer session.Close()

	if err == nil {
		res, err := r.Table("samples").Filter(map[string]interface{}{
			"file": map[string]interface{}{
				hashType: hash,
			},
		}).Run(session)
		utils.Assert(err)
		defer res.Close()

		// Scan query result into the person variable
		var samples []interface{}
		err = res.All(&samples)
		if err != nil {
			log.Errorf("Error scanning database result: %s\n", err)
			return r.WriteResponse{}
		}

		log.Debugf("%d samples\n", len(samples))
		log.Debugln("res: ", res)

		if res.IsNil() {
			// upsert into RethinkDB
			resp, err := r.Table("samples").Insert(
				map[string]interface{}{
					// "id":      sample.SHA256,
					"file": map[string]interface{}{
						hashType: hash,
					},
					"plugins": database.GetPluginsByCategory(),
				}, r.InsertOpts{Conflict: "replace"}).RunWrite(session)
			utils.Assert(err)

			return resp
		}
		log.Debugf("Hash %s already exists in the database.", hash)
		return r.WriteResponse{}

	}
	log.Debug(err)

	return r.WriteResponse{}
}
