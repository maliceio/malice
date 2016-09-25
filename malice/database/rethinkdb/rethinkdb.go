package rethinkdb

import (
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/go-plugin-utils/utils"
	"github.com/maliceio/malice/malice/database"
	"github.com/maliceio/malice/malice/persist"
	"github.com/maliceio/malice/utils"
	r "gopkg.in/dancannon/gorethink.v2"
)

// RethinkAddr RethinkDB address to user for connections
var RethinkAddr string

func init() {
	r.Log.Out = ioutil.Discard
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

	hashType, err := util.GetHashType(hash)
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
