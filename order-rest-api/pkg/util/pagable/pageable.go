package pagable

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	defaultSize = 10
	maxSize     = 100
	defaultPage = 1
)

type PageableQuery struct {
	Page string `json:"page"`
	Size string `json:"size"`
}

type Query struct {
	Page              int         `json:"page"`
	Size              int         `json:"size"`
	ExpressionFilters []Filter    `json:"filters"`
	ormConditions     interface{} `json:"-"`
}

type ListResponse struct {
	Items   interface{} `json:"items"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
	HasMore bool        `json:"has_more"`
}

// SetSize Set page size
func (q *Query) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}

	n, err := strconv.ParseUint(sizeQuery, 10, 32)
	if err != nil {
		return err
	}

	q.Size = int(n)
	if q.Size > maxSize {
		q.Size = maxSize
	}

	return nil
}

// SetPage Set page number
func (q *Query) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Page = defaultPage
		return nil
	}
	n, err := strconv.ParseUint(pageQuery, 10, 32)
	if err != nil {
		return err
	}
	q.Page = int(n)

	return nil
}

// GetOffset Get offset
func (q *Query) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// GetLimit Get limit
func (q *Query) GetLimit() int {
	return q.Size
}

// GetPage Get OrderBy
func (q *Query) GetPage() int {
	return q.Page
}

// Get OrderBy
func (q *Query) GetSize() int {
	if q.Size == 0 {
		return defaultSize
	}
	return q.Size
}

// Get total pages int
func (q *Query) GetTotalPages(totalCount int) int {
	d := float64(totalCount) / float64(q.GetSize())
	return int(math.Ceil(d))
}

// Get has more
func (q *Query) GetHasMore(total int) bool {
	return q.Page < total/q.GetSize()
}

func (q *Query) ORMConditions() interface{} {
	if q.ormConditions != nil {
		return q.ormConditions
	}
	var conditions []string
	for _, filter := range q.ExpressionFilters {
		var value interface{}

		// Type assertion to convert interface{} to []string
		if strSlice, ok := filter.Value.([]string); ok {
			value = ArrayToString(strSlice)
		} else {
			if filter.Field == "status" {
				value, _ = strconv.Atoi(fmt.Sprintf("%v", filter.Value))
			} else {
				value = fmt.Sprintf("'%s'", filter.Value)
			}
		}

		condition := q.ParseCondition(filter, value)
		conditions = append(conditions, condition)
	}

	q.ormConditions = strings.Join(conditions, " AND ")
	return q.ormConditions
}

func (q *Query) ParseCondition(ft Filter, value interface{}) string {
	condition := ""

	switch ft.Operation {
	case Equal:
		condition = ft.Field + " = " + fmt.Sprintf("%v", value)
	case NotEqual:
		condition = ft.Field + " <> " + fmt.Sprintf("%s", value)
	case LT:
		condition = ft.Field + " < " + fmt.Sprintf("%s", value)
	case LTE:
		condition = ft.Field + " <= " + fmt.Sprintf("%s", value)
	case GT:
		condition = ft.Field + " > " + fmt.Sprintf("%s", value)
	case GTE:
		condition = ft.Field + " >= " + fmt.Sprintf("%s", value)
	case In:
		condition = ft.Field + " IN " + fmt.Sprintf("(%v)", value)
	case NotIn:
		condition = ft.Field + " NOT IN " + fmt.Sprintf("%v", value)
	case Contains:
		condition = ft.Field + " LIKE " + "%" + fmt.Sprintf("%v", value) + "%"
	case NotContains:
		condition = ft.Field + " NOT LIKE " + "%" + fmt.Sprintf("%v", value) + "%"
	case IsNull:
		condition = ft.Field + " IS NULL"
	case IsNotNull:
		condition = ft.Field + " IS NOT NULL"
	case StartsWith:
		condition = ft.Field + " LIKE " + fmt.Sprintf(`'%s%s'`, ft.Value, "%")
	case EndsWith:
		condition = ft.Field + " LIKE " + fmt.Sprintf("'%s%s'", "%", ft.Value)
	}
	return condition
}

func (q *Query) UserORMConditions() interface{} {
	if q.ormConditions != nil {
		return q.ormConditions
	}
	var conditions []string
	for _, filter := range q.ExpressionFilters {

		if q.isUserRequest(filter.Field) {
			var value interface{}

			// Type assertion to convert interface{} to []string
			if strSlice, ok := filter.Value.([]string); ok {
				value = ArrayToString(strSlice)
			} else {
				if filter.Field == "status" {
					value, _ = strconv.Atoi(fmt.Sprintf("%v", filter.Value))
				} else {
					value = fmt.Sprintf("'%s'", filter.Value)
				}
			}

			condition := q.ParseCondition(filter, value)
			conditions = append(conditions, condition)
		}

	}

	q.ormConditions = strings.Join(conditions, " AND ")
	return q.ormConditions
}

func (q *Query) isUserRequest(fieldName string) bool {
	if fieldName == "status" || fieldName == "created_at" ||
		fieldName == "keyword" || fieldName == "order_uuid" || fieldName == "payment_method" {
		return true
	}
	return false
}

func (q *Query) ParseQueryParams() (map[string]string, error) {
	conditions := map[string]string{}
	for _, filter := range q.ExpressionFilters {
		if filter.Operation != Equal {
			return conditions, errors.New("only $eq filtering is supported")
		}
		conditions[filter.Field] = fmt.Sprint(filter.Value)
	}
	return conditions, nil
}
