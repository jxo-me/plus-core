package filestore

import (
	"context"
	"github.com/jxo-me/plus-core/pkg/tus"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test interface implementation of Filestore
var _ tus.DataStore = FileStore{}
var _ tus.TerminaterDataStore = FileStore{}
var _ tus.ConcaterDataStore = FileStore{}
var _ tus.LengthDeferrerDataStore = FileStore{}

func TestFilestore(t *testing.T) {
	a := assert.New(t)

	tmp, err := os.MkdirTemp("", "tus-filestore-")
	a.NoError(err)

	store := FileStore{tmp}
	ctx := context.Background()

	// Create new upload
	upload, err := store.NewUpload(ctx, tus.FileInfo{
		Size: 42,
		MetaData: map[string]string{
			"hello": "world",
		},
	})
	a.NoError(err)
	a.NotEqual(nil, upload)

	// Check info without writing
	info, err := upload.GetInfo(ctx)
	a.NoError(err)
	a.EqualValues(42, info.Size)
	a.EqualValues(0, info.Offset)
	a.Equal(tus.MetaData{"hello": "world"}, info.MetaData)
	a.Equal(2, len(info.Storage))
	a.Equal("filestore", info.Storage["Type"])
	a.Equal(filepath.Join(tmp, info.ID), info.Storage["Path"])

	// Write data to upload
	bytesWritten, err := upload.WriteChunk(ctx, 0, strings.NewReader("hello world"))
	a.NoError(err)
	a.EqualValues(len("hello world"), bytesWritten)

	// Check new offset
	info, err = upload.GetInfo(ctx)
	a.NoError(err)
	a.EqualValues(42, info.Size)
	a.EqualValues(11, info.Offset)

	// Read content
	reader, err := upload.GetReader(ctx)
	a.NoError(err)

	content, err := io.ReadAll(reader)
	a.NoError(err)
	a.Equal("hello world", string(content))
	reader.(io.Closer).Close()

	// Terminate upload
	a.NoError(store.AsTerminatableUpload(upload).Terminate(ctx))

	// Test if upload is deleted
	upload, err = store.GetUpload(ctx, info.ID)
	a.Equal(nil, upload)
	a.Equal(tus.ErrNotFound, err)
}

func TestMissingPath(t *testing.T) {
	a := assert.New(t)

	store := FileStore{"./path-that-does-not-exist"}
	ctx := context.Background()

	upload, err := store.NewUpload(ctx, tus.FileInfo{})
	a.Error(err)
	a.Equal("upload directory does not exist: ./path-that-does-not-exist", err.Error())
	a.Equal(nil, upload)
}

func TestNotFound(t *testing.T) {
	a := assert.New(t)

	store := FileStore{"./path"}
	ctx := context.Background()

	upload, err := store.GetUpload(ctx, "upload-that-does-not-exist")
	a.Error(err)
	a.Equal(tus.ErrNotFound, err)
	a.Equal(nil, upload)
}

func TestConcatUploads(t *testing.T) {
	a := assert.New(t)

	tmp, err := os.MkdirTemp("", "tus-filestore-concat-")
	a.NoError(err)

	store := FileStore{tmp}
	ctx := context.Background()

	// Create new upload to hold concatenated upload
	finUpload, err := store.NewUpload(ctx, tus.FileInfo{Size: 9})
	a.NoError(err)
	a.NotEqual(nil, finUpload)

	finInfo, err := finUpload.GetInfo(ctx)
	a.NoError(err)
	finId := finInfo.ID

	// Create three uploads for concatenating
	partialUploads := make([]tus.Upload, 3)
	contents := []string{
		"abc",
		"def",
		"ghi",
	}
	for i := 0; i < 3; i++ {
		upload, err := store.NewUpload(ctx, tus.FileInfo{Size: 3})
		a.NoError(err)

		n, err := upload.WriteChunk(ctx, 0, strings.NewReader(contents[i]))
		a.NoError(err)
		a.EqualValues(3, n)

		partialUploads[i] = upload
	}

	err = store.AsConcatableUpload(finUpload).ConcatUploads(ctx, partialUploads)
	a.NoError(err)

	// Check offset
	finUpload, err = store.GetUpload(ctx, finId)
	a.NoError(err)

	info, err := finUpload.GetInfo(ctx)
	a.NoError(err)
	a.EqualValues(9, info.Size)
	a.EqualValues(9, info.Offset)

	// Read content
	reader, err := finUpload.GetReader(ctx)
	a.NoError(err)

	content, err := io.ReadAll(reader)
	a.NoError(err)
	a.Equal("abcdefghi", string(content))
	reader.(io.Closer).Close()
}

func TestDeclareLength(t *testing.T) {
	a := assert.New(t)

	tmp, err := os.MkdirTemp("", "tus-filestore-declare-length-")
	a.NoError(err)

	store := FileStore{tmp}
	ctx := context.Background()

	upload, err := store.NewUpload(ctx, tus.FileInfo{
		Size:           0,
		SizeIsDeferred: true,
	})
	a.NoError(err)
	a.NotEqual(nil, upload)

	info, err := upload.GetInfo(ctx)
	a.NoError(err)
	a.EqualValues(0, info.Size)
	a.Equal(true, info.SizeIsDeferred)

	err = store.AsLengthDeclarableUpload(upload).DeclareLength(ctx, 100)
	a.NoError(err)

	updatedInfo, err := upload.GetInfo(ctx)
	a.NoError(err)
	a.EqualValues(100, updatedInfo.Size)
	a.Equal(false, updatedInfo.SizeIsDeferred)
}
