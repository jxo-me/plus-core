package cursor

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"time"
)

type MysqlCursorRepo struct {
	GDB gdb.DB
}

func NewMysqlCursorRepo(db gdb.DB) *MysqlCursorRepo {
	return &MysqlCursorRepo{GDB: db}
}

func (r *MysqlCursorRepo) GetLastPullTime(ctx context.Context, vendor string) (time.Time, error) {
	var pullTime time.Time
	result, err := r.GDB.Query(ctx, `SELECT pull_time FROM pull_cursor WHERE vendor = ? LIMIT 1`, vendor)
	if err != nil {
		return time.Now().UTC().Add(-48 * time.Hour), err
	}
	err = result.Structs(&pullTime)
	if err != nil {
		return time.Now().UTC().Add(-48 * time.Hour), err
	}

	return pullTime, err
}

func (r *MysqlCursorRepo) UpdatePullTime(ctx context.Context, vendor string, pullTime time.Time) error {
	_, err := r.GDB.Exec(ctx, `INSERT INTO pull_cursor (vendor, pull_time) VALUES (?, ?) ON DUPLICATE KEY UPDATE pull_time = VALUES(pull_time)`, vendor, pullTime)

	return err
}
