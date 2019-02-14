package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	docontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/puppetlabs/wash/datastore"
	"github.com/puppetlabs/wash/plugin"
	log "github.com/sirupsen/logrus"
)

type volume struct {
	plugin.EntryBase
	client    *client.Client
	startTime time.Time
	meta      datastore.Var
	list      datastore.Var
}

const mountpoint = "/mnt"

func newVolume(c *client.Client, v *types.Volume) (*volume, error) {
	startTime, err := time.Parse(time.RFC3339, v.CreatedAt)
	if err != nil {
		return nil, err
	}
	vol := &volume{
		EntryBase: plugin.NewEntry(v.Name),
		client:    c,
		startTime: startTime,
		meta:      datastore.NewVar(5 * time.Second),
		list:      datastore.NewVar(30 * time.Second),
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var metadata map[string]interface{}
	if err := json.Unmarshal(raw, &metadata); err != nil {
		return nil, err
	}
	vol.meta.Set(metadata)
	return vol, nil
}

func (v *volume) Metadata(ctx context.Context) (map[string]interface{}, error) {
	val, err := v.meta.Update(func() (interface{}, error) {
		_, raw, err := v.client.VolumeInspectWithRaw(ctx, v.Name())
		if err != nil {
			return nil, err
		}
		var metadata map[string]interface{}
		if err := json.Unmarshal(raw, &metadata); err != nil {
			return nil, err
		}
		return metadata, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(map[string]interface{}), nil
}

func (v *volume) Attr() plugin.Attributes {
	return plugin.Attributes{
		Ctime: v.startTime,
		Mtime: v.startTime,
		Atime: v.startTime,
	}
}

func (v *volume) LS(ctx context.Context) ([]plugin.Entry, error) {
	data, err := v.list.Update(func() (interface{}, error) {
		// Create a container that mounts a volume and inspects it. Run it and capture the output.
		cid, err := v.createContainer(ctx, plugin.StatCmd(mountpoint))
		if err != nil {
			return nil, err
		}
		defer v.client.ContainerRemove(ctx, cid, types.ContainerRemoveOptions{})

		log.Debugf("Starting container %v", cid)
		if err := v.client.ContainerStart(ctx, cid, types.ContainerStartOptions{}); err != nil {
			return nil, err
		}

		log.Debugf("Waiting for container %v", cid)
		waitC, errC := v.client.ContainerWait(ctx, cid, docontainer.WaitConditionNotRunning)
		var statusCode int64
		select {
		case err := <-errC:
			return nil, err
		case result := <-waitC:
			statusCode = result.StatusCode
			log.Debugf("Container %v finished[%v]: %v", cid, result.StatusCode, result.Error)
		}

		opts := types.ContainerLogsOptions{ShowStdout: true}
		if statusCode != 0 {
			opts.ShowStderr = true
		}

		log.Debugf("Gathering logs for %v", cid)
		output, err := v.client.ContainerLogs(ctx, cid, opts)
		if err != nil {
			return nil, err
		}
		defer output.Close()

		if statusCode != 0 {
			bytes, err := ioutil.ReadAll(output)
			if err != nil {
				return nil, err
			}
			return nil, errors.New(string(bytes))
		}

		return plugin.StatParseAll(output, mountpoint)
	})
	if err != nil {
		return nil, err
	}
	dirs := data.(plugin.DirMap)

	root := dirs[""]
	entries := make([]plugin.Entry, 0, len(root))
	for name, attr := range root {
		if attr.Mode.IsDir() {
			entries = append(entries, plugin.NewVolumeDir(name, attr, v.getContentCB(), "/"+name, dirs))
		} else {
			entries = append(entries, plugin.NewVolumeFile(name, attr, v.getContentCB(), "/"+name))
		}
	}
	return entries, nil
}

// Create a container that mounts a volume to a default mountpoint and runs a command.
func (v *volume) createContainer(ctx context.Context, cmd []string) (string, error) {
	// Use tty to avoid messing with the extra log formatting.
	cfg := docontainer.Config{Image: "busybox", Cmd: cmd, Tty: true}
	mounts := []mount.Mount{mount.Mount{
		Type:     mount.TypeVolume,
		Source:   v.Name(),
		Target:   mountpoint,
		ReadOnly: true,
	}}
	hostcfg := docontainer.HostConfig{Mounts: mounts}
	netcfg := network.NetworkingConfig{}
	created, err := v.client.ContainerCreate(ctx, &cfg, &hostcfg, &netcfg, "")
	if err != nil {
		return "", err
	}
	for _, warn := range created.Warnings {
		log.Debugf("Warning creating %v: %v", created.ID, warn)
	}
	return created.ID, nil
}

func (v *volume) getContentCB() plugin.ContentCB {
	return func(ctx context.Context, path string) (plugin.SizedReader, error) {
		// Create a container that mounts a volume and waits. Use it to download a file.
		cid, err := v.createContainer(ctx, []string{"sleep", "60"})
		if err != nil {
			return nil, err
		}
		defer v.client.ContainerRemove(ctx, cid, types.ContainerRemoveOptions{})

		log.Debugf("Starting container %v", cid)
		if err := v.client.ContainerStart(ctx, cid, types.ContainerStartOptions{}); err != nil {
			return nil, err
		}
		defer v.client.ContainerKill(ctx, cid, "SIGKILL")

		// Download file, then kill container.
		rdr, _, err := v.client.CopyFromContainer(ctx, cid, mountpoint+path)
		if err != nil {
			return nil, err
		}
		defer rdr.Close()

		tarReader := tar.NewReader(rdr)
		// Only expect one file.
		if _, err := tarReader.Next(); err != nil {
			return nil, err
		}
		bits, err := ioutil.ReadAll(tarReader)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(bits), nil
	}
}
