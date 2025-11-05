package sqlk

type FilterParser = func(queryValue string) (WhereWriter, []interface{})
