package sqlk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type DatabaseContext struct {
	conn *sqlx.DB
	ctx  context.Context
}

// MustPrepare prepare sql statements or exit app if fails or error
func (s *DatabaseContext) MustPrepare(query string) *sqlx.Stmt {
	stmt, err := s.conn.PreparexContext(s.ctx, query)
	if err != nil {
		panic(fmt.Errorf("sqlk: error while preparing statment [%s] (%s)", query, err))
	}
	return stmt
}

// MustPrepareFmt prepare sql statements from string format or exit app if fails or error
func (s *DatabaseContext) MustPrepareFmt(queryFmt string, args ...interface{}) *sqlx.Stmt {
	query := fmt.Sprintf(queryFmt, args...)
	return s.MustPrepare(query)
}

// MustPrepareNamedFmt prepare sql statements from string format with named bindvars or exit app if fails or error
func (s *DatabaseContext) MustPrepareNamedFmt(queryFmt string, args ...interface{}) *sqlx.NamedStmt {
	query := fmt.Sprintf(queryFmt, args...)
	return s.MustPrepareNamed(query)
}

// MustPrepareNamed prepare sql statements with named bindvars or exit app if fails or error
func (s *DatabaseContext) MustPrepareNamed(query string) *sqlx.NamedStmt {
	stmt, err := s.conn.PrepareNamedContext(s.ctx, query)
	if err != nil {
		panic(fmt.Errorf("sqlk: error while preparing named statment [%s] (%s)", query, err))
	}
	return stmt
}

// MustPrepareReplace prepare sql statements from a string replacement or exit app if fails or error
func (s *DatabaseContext) MustPrepareReplace(q string, values map[string]string) *sqlx.Stmt {
	for a, v := range values {
		q = strings.ReplaceAll(q, ":"+a, v)
	}
	return s.MustPrepare(q)
}

// MustPrepareRebind prepare sql statements and rebind with database adapter or exit app if fails or error
func (s *DatabaseContext) MustPrepareRebind(query string) *sqlx.Stmt {
	query = s.conn.Rebind(query)
	stmt, err := s.conn.PreparexContext(s.ctx, query)
	if err != nil {
		panic(fmt.Errorf("sqlk: error while prepare and rebind statment [%s] (%s)", query, err))
	}
	return stmt
}
