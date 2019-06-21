package aws

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/puppetlabs/wash/plugin"
)

// ec2InstanceMetadataJSON represents an EC2 instance's
// metadata.json file
type ec2InstanceMetadataJSON struct {
	plugin.EntryBase
	inst *ec2Instance
}

func newEC2InstanceMetadataJSON(ctx context.Context, inst *ec2Instance) (*ec2InstanceMetadataJSON, error) {
	im := &ec2InstanceMetadataJSON{
		EntryBase: plugin.NewEntry("metadata.json"),
	}
	im.DisableDefaultCaching()
	im.inst = inst

	content, err := im.Open(ctx)
	if err != nil {
		return nil, err
	}

	im.Attributes().SetSize(uint64(content.Size()))
	return im, nil
}

func (im *ec2InstanceMetadataJSON) Schema() *plugin.EntrySchema {
	return plugin.NewEntrySchema(im, "metadata.json").IsSingleton()
}

func (im *ec2InstanceMetadataJSON) Open(ctx context.Context) (plugin.SizedReader, error) {
	meta, err := plugin.CachedMetadata(ctx, im.inst)
	if err != nil {
		return nil, err
	}

	prettyMeta, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(prettyMeta), nil
}
