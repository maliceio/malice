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
