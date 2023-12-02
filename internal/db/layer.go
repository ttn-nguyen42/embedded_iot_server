package custdb

import (
	"context"
	"database/sql"
	"errors"
	custcon "labs/htmx-blog/internal/concurrent"
	custerror "labs/htmx-blog/internal/error"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/panjf2000/ants/v2"
)

var layeredOnce sync.Once

var layeredDb *LayeredDb

type LayeredDb struct {
	sqldb *sqlx.DB
	pool  *ants.Pool
}

func NewLayeredDb() *LayeredDb {
	return &LayeredDb{
		sqldb: Db(),
		pool:  custcon.New(100),
	}
}

func LayeredInit() {
	layeredOnce.Do(func() {
		layeredDb = NewLayeredDb()
	})
}

func Layered() *LayeredDb {
	return layeredDb
}

func (db *LayeredDb) Select(ctx context.Context, selectBuilder sq.SelectBuilder, dest interface{}) error {
	queryStr, arguments, err := selectBuilder.ToSql()
	if err != nil {
		return custerror.FormatInternalError("LayeredDb.Select: ToSql() err = %s", err)
	}

	sqldb := db.sqldb

	err = sqldb.SelectContext(ctx, dest, queryStr, arguments...)
	if err != nil {
		return custerror.FormatInternalError("LayeredDb.Select: SelectContext err = %s", err)
	}

	return nil
}

func (db *LayeredDb) Insert(ctx context.Context, insertBuilder sq.InsertBuilder) error {
	queryStr, arguments, err := insertBuilder.ToSql()
	if err != nil {
		return custerror.FormatInternalError("LayeredDb.Insert: ToSql() err = %s", err)
	}

	sqldb := db.sqldb

	_, err = sqldb.ExecContext(ctx, queryStr, arguments...)
	if err != nil {
		return custerror.FormatInternalError("LayeredDb.Insert: Exec() err = %s", err)
	}

	return nil
}

func (db *LayeredDb) Update(ctx context.Context, updateBuilder sq.UpdateBuilder) error {
	queryStr, arguments, err := updateBuilder.ToSql()
	if err != nil {
		return custerror.FormatInternalError("LayeredDb.Update: ToSql() err = %s", err)
	}

	sqldb := db.sqldb

	res, err := sqldb.ExecContext(ctx, queryStr, arguments...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return custerror.ErrorNotFound
		}
		return custerror.FormatInternalError("LayeredDb.Update: Exec() err = %s", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return custerror.FormatInternalError("LayeredDb.Update: RowsAffected() err = %s", err)
	}

	if affected == 0 {
		return custerror.ErrorNotFound
	}

	return nil
}

func (db *LayeredDb) Delete(ctx context.Context, deleteBuilder sq.DeleteBuilder) error {
	queryStr, arguments, err := deleteBuilder.ToSql()
	if err != nil {
		return custerror.FormatInternalError("LayeredDb.Delete: ToSql() err = %s", err)
	}

	sqldb := db.sqldb

	res, err := sqldb.ExecContext(ctx, queryStr, arguments...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return custerror.ErrorNotFound
		}
		return custerror.FormatInternalError("LayeredDb.Delete: Exec() err = %s", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return custerror.FormatInternalError("LayeredDb.Delete: RowsAffected() err = %s", err)
	}

	if affected == 0 {
		return custerror.ErrorNotFound
	}

	return nil
}
