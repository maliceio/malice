package search

import (
	"net/http"

	"github.com/maliceio/engine/api/server/httputils"
	"github.com/maliceio/engine/api/types/search"
	"golang.org/x/net/context"
)

func (sr *searchRouter) doSearch(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	config := &search.Config{
		Hash: r.FormValue("hash"),
	}
	result, err := sr.backend.Search(r.FormValue("path"), config)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, result)
}
