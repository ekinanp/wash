package api

import (
	"net/http"

	"github.com/puppetlabs/wash/activity"
	"github.com/puppetlabs/wash/plugin"
)

// swagger:route DELETE /fs/delete delete
//
// Deletes the entry at the specified path. If the entry is a Parent
// and the recurse option is not set, then Delete will return an error
// response.
//
//     Schemes: http
//
//     Responses:
//       200
//       400: errorResp
//       404: errorResp
//       500: errorResp
var deleteHandler handler = func(w http.ResponseWriter, r *http.Request) *errorResponse {
	ctx := r.Context()
	entry, path, errResp := getEntryFromRequest(r)
	if errResp != nil {
		return errResp
	}
	if !plugin.DeleteAction().IsSupportedOn(entry) {
		return unsupportedActionResponse(path, plugin.DeleteAction())
	}
	if plugin.ListAction().IsSupportedOn(entry) {
		// Parent entries can only be deleted if the recurse option is set.
		recurse, err := getBoolParam(r.URL, "recurse")
		if err != nil {
			return err
		}
		if !recurse {
			return badActionRequestResponse(
				path,
				plugin.DeleteAction(),
				"parent entries can only be deleted when the recurse option is set",
			)
		}
	}
	err := plugin.Delete(ctx, entry.(plugin.Deletable))
	if err != nil {
		return erroredActionResponse(path, plugin.DeleteAction(), err.Error())
	}
	activity.Record(ctx, "API: Delete %v", path)
	return nil
}
