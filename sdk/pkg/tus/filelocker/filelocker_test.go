package filelocker

import (
	"github.com/jxo-me/plus-core/sdk/pkg/tus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var _ tus.Locker = &FileLocker{}

func TestFileLocker(t *testing.T) {
	a := assert.New(t)

	dir, err := ioutil.TempDir("", "tus-file-locker")
	a.NoError(err)

	locker := FileLocker{dir}

	lock1, err := locker.NewLock("one")
	a.NoError(err)

	a.NoError(lock1.Lock())
	a.Equal(tus.ErrFileLocked, lock1.Lock())

	lock2, err := locker.NewLock("one")
	a.NoError(err)
	a.Equal(tus.ErrFileLocked, lock2.Lock())

	a.NoError(lock1.Unlock())
}
