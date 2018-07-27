package rethinkdb

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/go-plugin-utils/utils"
	r "gopkg.in/dancannon/gorethink.v2"
)

var (
	name     string
	category string
)

type PluginResults struct {
	ID   string `json:"id" gorethink:"id,omitempty"`
	Data map[string]interface{}
}

// WritePluginResultsToDatabase upserts plugin results into Database
func WritePluginResultsToDatabase(results map[string]interface{}) {

	// connect to RethinkDB
	session, err := r.Connect(r.ConnectOpts{
		Address:  fmt.Sprintf("%s:28015", utils.Getopt("MALICE_RETHINKDB", "rethink")),
		Timeout:  5 * time.Second,
		Database: "malice",
	})
	if err != nil {
		log.Debug(err)
		return
	}
	defer session.Close()

	res, err := r.Table("samples").Get(results["ID"]).Run(session)
	utils.Assert(err)
	defer res.Close()

	if res.IsNil() {
		// upsert into RethinkDB
		resp, err := r.Table("samples").Insert(results, r.InsertOpts{Conflict: "replace"}).RunWrite(session)
		utils.Assert(err)
		log.Debug(resp)
	} else {
		resp, err := r.Table("samples").Get(results["ID"]).Update(map[string]interface{}{
			"plugins": map[string]interface{}{
				category: map[string]interface{}{
					name: results["Data"],
				},
			},
		}).RunWrite(session)
		utils.Assert(err)

		log.Debug(resp)
	}
}
