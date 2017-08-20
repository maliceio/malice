package scan

import (
	"net/http"

	"github.com/maliceio/engine/api/server/httputils"
	"github.com/maliceio/engine/api/types/scan"
	"golang.org/x/net/context"
)

func (sr *scanRouter) doScan(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	config := &scan.Config{
		Path: r.FormValue("path"),
	}
	result, err := sr.backend.Scan(r.FormValue("path"), config)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, result)
}
