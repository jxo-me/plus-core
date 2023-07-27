package config

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/plus-core/pkg/v2/tus"
	"github.com/jxo-me/plus-core/pkg/v2/tus/filestore"
)

const (
	UploadCfgName = "uploadConfig"
)

var insUpload = Upload{}

type Upload struct{}

func UploadConfig() *Upload {
	return &insUpload
}

func (q *Upload) String() string {
	return UploadCfgName
}

func (q *Upload) Init(ctx context.Context) error {
	cf := tus.Config{}
	c, err := Setting().Cfg().Get(ctx, "settings.uploads.tus.default", "")
	if err != nil {
		return err
	}
	err = c.Scan(&cf)
	if err != nil {
		return err
	}
	// Create a new FileStore instance which is responsible for
	// storing the uploaded file on disk in the specified directory.
	// This path _must_ exist before tusd will store uploads in it.
	// If you want to save them on a different medium, for example
	// a remote FTP server, you can implement your own storage backend
	// by implementing the tusd.DataStore interface.
	store := filestore.FileStore{
		Path: cf.Path,
	}

	// A storage backend for tusd may consist of multiple different parts which
	// handle upload creation, locking, termination, and so on. The composer is a
	// place where all those separated pieces are joined together. In this example,
	// we only use the file store, but you may plug in multiple.
	composer := tus.NewStoreComposer()
	store.UseIn(composer)
	cf.StoreComposer = composer

	logger := glog.New()
	err = logger.SetConfigWithMap(g.Map{
		"path":   cf.LogPath,
		"file":   cf.LogFile,
		"level":  cf.LogLevel,
		"stdout": cf.LogStdout,
	})
	if err != nil {
		return err
	}
	cf.Logger = logger
	Setting().Config().Tus = cf
	return nil
}
