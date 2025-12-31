package query

import "github.com/konsultin/project-goes-here/libs/sqlk/option"

func On(col string, args ...interface{}) option.SetOptionFn {
	return func(o *option.Options) {
		// Evaluate options
		opts := option.EvaluateOptions(args)
		schema := opts.GetSchema()

		// Get table name
		var tableName string
		if schema == nil {
			tableName = joinTableFlag
		}

		// Set variable
		o.KV[option.VariableKey] = &columnWriter{
			name:      col,
			tableName: tableName,
		}
	}
}
