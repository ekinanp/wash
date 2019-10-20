package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/puppetlabs/wash/activity"
	apitypes "github.com/puppetlabs/wash/api/types"
	"github.com/puppetlabs/wash/plugin"
)

// swagger:response
//nolint:deadcode,unused
type entryList struct {
	// in: body
	Entries []apitypes.Entry
}

// swagger:route GET /fs/list list listEntries
//
// Lists children of a path
//
// Returns a list of Entry objects describing children of the given path.
//
//     Produces:
//     - application/json
//
//     Schemes: http
//
//     Responses:
//       200: entryList
//       400: errorResp
//       404: errorResp
//       500: errorResp
var listHandler handler = func(w http.ResponseWriter, r *http.Request) *errorResponse {
	ctx := r.Context()
	entry, path, errResp := getEntryFromRequest(r)
	if errResp != nil {
		return errResp
	}

	if !plugin.ListAction().IsSupportedOn(entry) {
		return unsupportedActionResponse(path, plugin.ListAction())
	}

	parent := entry.(plugin.Parent)
	entries, err := plugin.ListWithAnalytics(ctx, parent)
	if err != nil {
		if cnameErr, ok := err.(plugin.DuplicateCNameErr); ok {
			return duplicateCNameResponse(cnameErr)
		}

		return erroredActionResponse(path, plugin.ListAction(), err.Error())
	}

	result := make([]apitypes.Entry, 0, entries.Len())
	entries.Range(func(_ string, entry plugin.Entry) bool {
		apiEntry := toAPIEntry(entry)
		apiEntry.Path = path + "/" + apiEntry.CName
		result = append(result, apiEntry)
		return true
	})
	// Sort entries so they have a deterministic order.
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	activity.Record(ctx, "API: List %v %+v", path, result)

	jsonEncoder := json.NewEncoder(w)
	if err = jsonEncoder.Encode(result); err != nil {
		return unknownErrorResponse(fmt.Errorf("Could not marshal list results for %v: %v", path, err))
	}
	return nil
}
