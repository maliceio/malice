package database

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/malice/persist"
	"github.com/maliceio/malice/plugins"
	"github.com/maliceio/malice/utils"
	r "gopkg.in/dancannon/gorethink.v2"
)

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getPluginsByCategory() map[string]interface{} {
	categoryList := make(map[string]interface{})
	for _, category := range plugins.GetCategories() {
		pluginList := make(map[string]interface{})
		for _, plugin := range plugins.GetAllPluginsInCategory(category) {
			pluginList[plugin.Name] = nil
		}
		categoryList[category] = pluginList
	}

	return categoryList
}

func getPlugins() map[string]interface{} {
	pluginList := make(map[string]interface{})
	for _, plugin := range plugins.Plugs.Plugins {
		pluginList[plugin.Name] = nil
	}

	return pluginList
}

// TestConnection tests the rethinkDB connection
func TestConnection() error {
	// connect to RethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address: fmt.Sprintf("%s:28015", getopt("MALICE_RETHINKDB", "rethink")),
	})
	defer session.Close()

	return err
}

// InitRethinkDB initalizes rethinkDB for use with malice
func InitRethinkDB() error {
	// connect to RethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address: fmt.Sprintf("%s:28015", getopt("MALICE_RETHINKDB", "rethink")),
	})
	defer session.Close()
	assert(err)
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
	session, err := r.Connect(r.ConnectOpts{
		Address:  fmt.Sprintf("%s:28015", getopt("MALICE_RETHINKDB", "rethink")),
		Database: "malice",
	})
	defer session.Close()

	if err == nil {
		res, err := r.Table("samples").Filter(map[string]interface{}{
			"file": map[string]interface{}{
				"sha256": sample.SHA256,
			},
		}).Run(session)
		assert(err)
		defer res.Close()

		// Scan query result into the person variable
		var samples []interface{}
		err = res.All(&samples)
		if err != nil {
			fmt.Printf("Error scanning database result: %s\n", err)
			return r.WriteResponse{}
		}
		fmt.Printf("%d samples\n", len(samples))

		fmt.Println("res: ", res)

		if res.IsNil() {
			// upsert into RethinkDB
			resp, err := r.Table("samples").Insert(
				map[string]interface{}{
					// "id":      sample.SHA256,
					"file":    sample,
					"plugins": getPluginsByCategory(),
				}, r.InsertOpts{Conflict: "replace"}).RunWrite(session)
			assert(err)

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
	assert(err)

	// connect to RethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address:  fmt.Sprintf("%s:28015", getopt("MALICE_RETHINKDB", "rethink")),
		Database: "malice",
	})
	defer session.Close()

	if err == nil {
		res, err := r.Table("samples").Filter(map[string]interface{}{
			"file": map[string]interface{}{
				hashType: hash,
			},
		}).Run(session)
		assert(err)
		defer res.Close()

		// Scan query result into the person variable
		var samples []interface{}
		err = res.All(&samples)
		if err != nil {
			fmt.Printf("Error scanning database result: %s\n", err)
			return r.WriteResponse{}
		}
		fmt.Printf("%d samples\n", len(samples))

		fmt.Println("res: ", res)

		if res.IsNil() {
			// upsert into RethinkDB
			resp, err := r.Table("samples").Insert(
				map[string]interface{}{
					// "id":      sample.SHA256,
					"file": map[string]interface{}{
						hashType: hash,
					},
					"plugins": getPluginsByCategory(),
				}, r.InsertOpts{Conflict: "replace"}).RunWrite(session)
			assert(err)

			return resp
		}
		log.Debugf("Hash %s already exists in the database.", hash)
		return r.WriteResponse{}

	}
	log.Debug(err)

	return r.WriteResponse{}
}
