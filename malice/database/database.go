package database

import (
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/malice/malice/persist"
	"github.com/maliceio/malice/plugins"
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

// InitRethinkDB initalizes rethinkDB for use with malice
func InitRethinkDB() error {
	// connect to RethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address: fmt.Sprintf("%s:28015", getopt("MALICE_RETHINKDB", "rethink")),
		Timeout: 5 * time.Second,
	})
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

// WriteToDatabase inserts sample into Database
func WriteToDatabase(sample persist.File) {

	// connect to RethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address:  fmt.Sprintf("%s:28015", getopt("MALICE_RETHINKDB", "rethink")),
		Timeout:  5 * time.Second,
		Database: "malice",
	})
	if err == nil {
		res, err := r.Table("samples").Get(sample.SHA256).Run(session)
		assert(err)
		defer res.Close()

		if res.IsNil() {
			// upsert into RethinkDB
			resp, err := r.Table("samples").Insert(
				map[string]interface{}{
					"id":      sample.SHA256,
					"file":    sample,
					"plugins": getPluginsByCategory(),
				}, r.InsertOpts{Conflict: "replace"}).RunWrite(session)
			assert(err)
			log.Debug(resp)
		} else {
			log.Debugf("Sample: %s already exists in the database.", sample.SHA256)
		}

	} else {
		log.Debug(err)
	}
}
