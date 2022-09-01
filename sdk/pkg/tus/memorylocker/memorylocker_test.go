package memorylocker

import (
	"github.com/jxo-me/plus-core/sdk/pkg/tus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var _ tus.Locker = &MemoryLocker{}

func TestMemoryLocker(t *testing.T) {
	a := assert.New(t)

	locker := New()

	lock1, err := locker.NewLock("one")
	a.NoError(err)

	a.NoError(lock1.Lock())
	a.Equal(tus.ErrFileLocked, lock1.Lock())

	lock2, err := locker.NewLock("one")
	a.NoError(err)
	a.Equal(tus.ErrFileLocked, lock2.Lock())

	a.NoError(lock1.Unlock())
	a.NoError(lock1.Unlock())
}
