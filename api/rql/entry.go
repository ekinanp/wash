package rql

import apitypes "github.com/puppetlabs/wash/api/types"

type EntrySchema = apitypes.EntrySchema

type Entry struct {
	apitypes.Entry
	Schema *EntrySchema
}
