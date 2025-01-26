package utilorm

import (
	"errors"
	"go-blog-admin/internal/util/utilpaging"
	"go-blog-admin/internal/util/utilunsafe"
	"strconv"
	"strings"
)

type FieldType int32
type FieldName string
type FieldOptions map[FieldName]FieldType
type SortOptions map[string]string

const (
	IntType FieldType = iota // iota starts from 0 and increments by 1
	StringType
	ShortStringType
	BoolType

	// short
	// float
	// date
)

func (x FieldName) String() string {
	return string(x)
}

func (x FieldType) Parse(str string) (any, error) {
	switch x {
	case IntType:
		return strconv.Atoi(str)
	case StringType:
		return str, nil
	case ShortStringType:
		return utilunsafe.UnsafeLen(str), nil
	case BoolType:
		return strconv.ParseBool(str)
	default:
		return nil, errors.New("undefined parser type")
	}
}
func LimitString(str string) string {
	return str
}

type SQLBuilder struct {
	WhereText strings.Builder
	WhereArgs []any
}

// WhereSearch create " and ((title ilike ?) or (content_markdown ilike ?))"
func WhereSearch(sql *SQLBuilder, search string,

	fieldOptions FieldOptions,

) (err error) {

	if search == "" {
		return nil
	}

	search = utilunsafe.UnsafeLen(search) // unsafe check

	if !strings.ContainsAny(search, "_%") {
		search = "%" + search + "%"
	}

	first := true
	changed := false

	for fieldName, fieldType := range fieldOptions {

		if fieldType != ShortStringType {
			continue
		}

		if first {
			sql.WhereText.WriteString(" AND (")
		}

		{
			changed = true

			if !first {
				sql.WhereText.WriteString(" OR ")
			}

			if err = utilunsafe.UnsafeName(fieldName.String()); err != nil {
				return err
			}

			sql.WhereText.WriteString("(" + fieldName.String() + " ILIKE ?)")
			sql.WhereArgs = append(sql.WhereArgs, search)
		}

		first = false

	}

	if changed {
		sql.WhereText.WriteString(")")
	}

	return err
}

// WhereFromFilter apply bool filter
func WhereFromFilter(
	sql *SQLBuilder,
	filter *utilpaging.PagingInputDTO,
	fieldOptions FieldOptions,
) (err error) {

	for fieldName, fieldType := range fieldOptions {
		var value any
		var exists bool
		value, exists, err = utilpaging.FilterParser(filter.Filters, "filter_"+fieldName.String(), fieldType.Parse)
		if exists && err == nil {

			if err = utilunsafe.UnsafeName(fieldName.String()); err != nil {
				return err
			}

			sql.WhereText.WriteString(" AND (" + fieldName.String() + " = ?)")
			sql.WhereArgs = append(sql.WhereArgs, value)
		}
		if err != nil {
			break
		}
	}
	return err
}
