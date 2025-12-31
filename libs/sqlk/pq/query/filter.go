package query

import (
	"fmt"

	"github.com/konsultin/project-goes-here/libs/sqlk"
	"github.com/konsultin/project-goes-here/libs/sqlk/op"
	"github.com/konsultin/project-goes-here/libs/sqlk/option"
	qs "github.com/konsultin/project-goes-here/libs/sqlk/parse/querystring"
	"github.com/konsultin/project-goes-here/libs/sqlk/schema"
)

type FilterBuilder struct {
	conditions []sqlk.WhereWriter
	args       []interface{}
}

func (b *FilterBuilder) Conditions() sqlk.WhereWriter {
	return And(b.conditions...)
}

func (b *FilterBuilder) Args() []interface{} {
	return b.args
}

// NewFilter create a FilterBuilder that convert querystring to WHERE conditions
func NewFilter(qs map[string]string, funcMap map[string]sqlk.FilterParser) *FilterBuilder {
	b := FilterBuilder{
		conditions: make([]sqlk.WhereWriter, 0),
		args:       make([]interface{}, 0),
	}

	// Get value from query string
	for k, v := range qs {
		// If value is empty string, then skip
		if v == "" {
			continue
		}

		// Get mapper function
		fn, ok := funcMap[k]

		// If function mapper is not set, then skip
		if !ok {
			continue
		}

		w, args := fn(v)

		// If writer is empty, then skip
		if w == nil {
			continue
		}

		// Append condition
		b.conditions = append(b.conditions, w)

		// If arguments is exists, then merge arguments
		if len(args) > 0 {
			b.args = append(b.args, args...)
		}
	}

	return &b
}

func LikeFilter(col string, likeVar op.LikeVariable, args ...interface{}) sqlk.FilterParser {
	// Get options
	opts := option.EvaluateOptions(args)
	s := opts.GetSchema()

	return func(qv string) (sqlk.WhereWriter, []interface{}) {
		// Trim value
		switch likeVar {
		case op.LikeSubString:
			qv = fmt.Sprintf(`%%%s%%`, qv)
		case op.LikePrefix:
			qv = fmt.Sprintf(`%%%s`, qv)
		case op.LikeSuffix:
			qv = fmt.Sprintf(`%s%%`, qv)
		}

		w := ILike(Column(col, option.Schema(s)), qv)

		return w, []interface{}{qv}
	}
}

func EqualFilter(s *schema.Schema, col string) sqlk.FilterParser {
	return func(qv string) (sqlk.WhereWriter, []interface{}) {
		w := Equal(Column(col, option.Schema(s)))
		return w, []interface{}{qv}
	}
}

func TimeGreaterThanEqualFilter(s *schema.Schema, col string, args ...string) sqlk.FilterParser {
	return func(qv string) (sqlk.WhereWriter, []interface{}) {
		// Parse time
		t, ok := qs.ParseTime(qv, args...)
		if !ok {
			return nil, nil
		}

		// Create schema
		w := GreaterThanEqual(Column(col, option.Schema(s)))
		return w, []interface{}{t.UTC()}
	}
}

func TimeLessThanEqualFilter(s *schema.Schema, col string, args ...string) sqlk.FilterParser {
	return func(qv string) (sqlk.WhereWriter, []interface{}) {
		// Parse time
		t, ok := qs.ParseTime(qv, args...)
		if !ok {
			return nil, nil
		}

		w := LessThanEqual(Column(col, option.Schema(s)))
		return w, []interface{}{t.UTC()}
	}
}

func IntGreaterThanEqualFilter(s *schema.Schema, col string) sqlk.FilterParser {
	return func(qv string) (sqlk.WhereWriter, []interface{}) {
		// Parse int value
		i, ok := qs.ParseInt(qv)
		if !ok {
			return nil, nil
		}

		w := GreaterThanEqual(Column(col, option.Schema(s)))
		return w, []interface{}{i}
	}
}

func IntLessThanEqualFilter(s *schema.Schema, col string) sqlk.FilterParser {
	return func(qv string) (sqlk.WhereWriter, []interface{}) {
		// Parse int value
		i, ok := qs.ParseInt(qv)
		if !ok {
			return nil, nil
		}

		w := LessThanEqual(Column(col, option.Schema(s)))
		return w, []interface{}{i}
	}
}

func FloatGreaterThanEqualFilter(s *schema.Schema, col string) sqlk.FilterParser {
	return func(qv string) (sqlk.WhereWriter, []interface{}) {
		// Parse float value
		f, ok := qs.ParseFloat(qv)
		if !ok {
			return nil, nil
		}

		w := GreaterThanEqual(Column(col, option.Schema(s)))
		return w, []interface{}{f}
	}
}

func FloatLessThanEqualFilter(s *schema.Schema, col string) sqlk.FilterParser {
	return func(qv string) (sqlk.WhereWriter, []interface{}) {
		// Parse float value
		f, ok := qs.ParseFloat(qv)
		if !ok {
			return nil, nil
		}

		w := LessThanEqual(Column(col, option.Schema(s)))
		return w, []interface{}{f}
	}
}
