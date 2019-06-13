package find

import (
	"github.com/puppetlabs/wash/api/client"
	"github.com/puppetlabs/wash/cmd/internal/find/types"
)

// info is a wrapper to c.Info
func info(c client.Client, path string) (types.Entry, error) {
	e, err := c.Info(path)
	if err != nil {
		return types.Entry{}, err
	}
	return types.NewEntry(e, path), nil
}

// list is a wrapper to c.List that handles normalizing the children's
// path relative to e's normalized path
func list(c client.Client, e types.Entry) ([]types.Entry, error) {
	rawChildren, err := c.List(e.Path)
	if err != nil {
		return nil, err
	}
	children := make([]types.Entry, len(rawChildren))
	for i, rawChild := range rawChildren {
		children[i] = types.NewEntry(rawChild, e.NormalizedPath+"/"+rawChild.CName)
	}
	return children, nil
}
