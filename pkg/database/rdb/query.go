package rdb

import (
	"errors"
	"fmt"
	"strings"
)

type entryPropsApplierFn[E any] func(E) []any

func BulkInsertQuery[E any](table string, columns []string, entries []E, propsApplier entryPropsApplierFn[E]) (string, []any, error) {
	if table == "" {
		return "", nil, errors.New("table name must be specified")
	}

	if len(columns) == 0 {
		return "", nil, errors.New("no columns provided")
	}

	if len(entries) == 0 {
		return "", nil, errors.New("no entries provided")
	}

	cols := buildInsertColumns(columns)
	placeholders, params := buildInsertParams(entries, propsApplier)

	q := fmt.Sprintf("INSERT INTO %s%s VALUES%s", table, cols, placeholders)
	return q, params, nil
}

func WhereIn[E any](inOpts []E) (string, []any, error) {
	if len(inOpts) == 0 {
		return "", nil, errors.New("no entries provided")
	}

	var q strings.Builder
	q.WriteString("(")

	params := make([]any, 0)
	lastIndex := len(inOpts) - 1
	for i, opt := range inOpts {
		params = append(params, opt)
		q.WriteString(fmt.Sprintf("$%d", i+1))

		if i != lastIndex {
			q.WriteString(",")
		}
	}
	q.WriteString(")")

	return q.String(), params, nil
}

func buildInsertColumns(columns []string) string {
	var cols strings.Builder
	lastIndex := len(columns) - 1

	cols.WriteString("(")
	for i, column := range columns {
		cols.WriteString(column)
		if i != lastIndex {
			cols.WriteString(",")
		}
	}
	cols.WriteString(")")

	return cols.String()
}

func buildInsertParams[E any](entries []E, propsApplier entryPropsApplierFn[E]) (string, []any) {
	var query strings.Builder
	params := make([]any, 0)
	lastIndex := len(entries) - 1

	for i, entry := range entries {
		props := propsApplier(entry)
		propsCount := len(props)
		position := i * propsCount

		var placeholders strings.Builder
		placeholders.WriteString("(")
		for j := 0; j < propsCount; j++ {
			placeholders.WriteString(fmt.Sprintf("$%d", position+j+1))
			if j != propsCount-1 {
				placeholders.WriteString(",")
			}
		}
		placeholders.WriteString(")")

		query.WriteString(placeholders.String())
		if i != lastIndex {
			query.WriteString(",")
		}
		params = append(params, props...)
	}

	return query.String(), params
}
