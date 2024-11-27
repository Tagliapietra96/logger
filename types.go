package logger

import (
	"fmt"
	"strings"
)

type QueryOperator interface {
	op() string
}

type NumericOperator int

const (
	EQUAL NumericOperator = iota
	NOT_EQUAL
	GREATER_THAN
	GREATER_THAN_OR_EQUAL
	LESS_THAN
	LESS_THAN_OR_EQUAL
)

func (no NumericOperator) op() string {
	var operator string

	switch no {
	case EQUAL:
		operator = "="
	case NOT_EQUAL:
		operator = "!="
	case GREATER_THAN:
		operator = ">"
	case GREATER_THAN_OR_EQUAL:
		operator = ">="
	case LESS_THAN:
		operator = "<"
	case LESS_THAN_OR_EQUAL:
		operator = "<="
	default:
		operator = ""
	}

	return operator
}

type StringOperator int

const (
	CONTAINS StringOperator = iota
	NOT_CONTAINS
	SAME
)

func (so StringOperator) op() string {
	var operator string

	switch so {
	case CONTAINS:
		operator = "LIKE"
	case NOT_CONTAINS:
		operator = "NOT LIKE"
	case SAME:
		operator = "="
	default:
		operator = ""
	}
	return operator
}

type SortOperator int

const (
	ASC SortOperator = iota
	DESC
)

func (so SortOperator) op() string {
	var operator string

	switch so {
	case ASC:
		operator = "ASC"
	case DESC:
		operator = "DESC"
	default:
		operator = ""
	}
	return operator
}

type QueryConfiguration func(*strings.Builder)

func CustomQuery(query string) QueryConfiguration {
	return func(sb *strings.Builder) {
		sb.WriteString(" ")
		sb.WriteString(query)
	}
}

func prepareFilter(config QueryConfiguration) QueryConfiguration {
	return func(sb *strings.Builder) {
		var filter, order, limit string
		s := sb.String()
		if s == "" {
			sb.WriteString(defaultQuery)
		}

		if strings.Contains(s, " WHERE ") {
			pieces := strings.Split(s, " WHERE ")
			filter = pieces[1]
			if strings.Contains(filter, " ORDER BY ") {
				pieces = strings.Split(filter, " ORDER BY ")
				filter = pieces[0]
				order = pieces[1]
				if strings.Contains(order, " LIMIT ") {
					pieces = strings.Split(order, " LIMIT ")
					order = pieces[0]
					limit = pieces[1]
				}
			} else if strings.Contains(filter, " LIMIT ") {
				pieces = strings.Split(filter, " LIMIT ")
				filter = pieces[0]
				limit = pieces[1]
			}
			sb.Reset()
			sb.WriteString(defaultQuery)
			sb.WriteString(" WHERE ")
			sb.WriteString(filter)
		} else if strings.Contains(s, " ORDER BY ") {
			pieces := strings.Split(s, " ORDER BY ")
			order = pieces[1]
			if strings.Contains(order, " LIMIT ") {
				pieces = strings.Split(order, " LIMIT ")
				order = pieces[0]
				limit = pieces[1]
			}
			sb.Reset()
			sb.WriteString(defaultQuery)
		} else if strings.Contains(s, " LIMIT ") {
			pieces := strings.Split(s, " LIMIT ")
			limit = pieces[1]
			sb.Reset()
			sb.WriteString(defaultQuery)
		}

		if !strings.Contains(s, " WHERE ") {
			sb.WriteString(" WHERE ")
		} else {
			sb.WriteString(" AND ")
		}

		config(sb)

		if order != "" {
			sb.WriteString(" ORDER BY ")
			sb.WriteString(order)
		}

		if limit != "" {
			sb.WriteString(" LIMIT ")
			sb.WriteString(limit)
		}
	}
}

func prepareSort(config QueryConfiguration) QueryConfiguration {
	return func(sb *strings.Builder) {
		var base, order, limit string
		s := sb.String()
		if s == "" {
			sb.WriteString(defaultQuery)
		}

		if strings.Contains(s, " ORDER BY ") {
			pieces := strings.Split(s, " ORDER BY ")
			base = pieces[0]
			order = pieces[1]
			if strings.Contains(order, " LIMIT ") {
				pieces = strings.Split(order, " LIMIT ")
				order = pieces[0]
				limit = pieces[1]
			}
			sb.Reset()
			sb.WriteString(base)
			sb.WriteString(" ORDER BY ")
			sb.WriteString(order)
		} else if strings.Contains(s, " LIMIT ") {
			pieces := strings.Split(s, " LIMIT ")
			base = pieces[0]
			limit = pieces[1]
			sb.Reset()
			sb.WriteString(base)
		}

		if !strings.Contains(s, " ORDER BY ") {
			sb.WriteString(" ORDER BY ")
		} else {
			sb.WriteString(", ")
		}

		config(sb)

		if limit != "" {
			sb.WriteString(" LIMIT ")
			sb.WriteString(limit)
		}
	}
}

func AddFilters(configs ...QueryConfiguration) QueryConfiguration {
	return func(sb *strings.Builder) {
		for _, config := range configs {
			prepareFilter(config)
		}
	}
}

func AddSorts(configs ...QueryConfiguration) QueryConfiguration {
	return func(sb *strings.Builder) {
		for _, config := range configs {
			prepareSort(config)
		}
	}
}

func AddLimit(limitAndOffset ...int) QueryConfiguration {
	return func(sb *strings.Builder) {
		if len(limitAndOffset) == 0 {
			return
		}

		var base string
		s := sb.String()
		if s == "" {
			sb.WriteString(defaultQuery)
		}

		if strings.Contains(s, " LIMIT ") {
			pieces := strings.Split(s, " LIMIT ")
			base = pieces[0]
			sb.Reset()
			sb.WriteString(base)
		}

		sb.WriteString(" LIMIT ")
		sb.WriteString(fmt.Sprintf("%d", limitAndOffset[0]))

		if len(limitAndOffset) > 1 {
			sb.WriteString(" OFFSET ")
			sb.WriteString(fmt.Sprintf("%d", limitAndOffset[1]))
		}
	}
}

func FilterByLevel(level LogLevel) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s = %d", Level.String(), level))
	})
}

func FilterByContext(context string, operator StringOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		if operator == CONTAINS || operator == NOT_CONTAINS {
			context = fmt.Sprintf("%%%s%%", context)
		}
		sb.WriteString(fmt.Sprintf("%s %s '%s'", Tags.String(), operator.op(), context))
	})
}

func FilterByCallerFile(file string, operator StringOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		if operator == CONTAINS || operator == NOT_CONTAINS {
			file = fmt.Sprintf("%%%s%%", file)
		}
		sb.WriteString(fmt.Sprintf("%s %s '%s'", CallerFile.String(), operator.op(), file))
	})
}

func FilterByCallerLine(line int, operator NumericOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s %d", Callerline.String(), operator.op(), line))
	})
}

func FilterByCallerFunction(function string, operator StringOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		if operator == CONTAINS || operator == NOT_CONTAINS {
			function = fmt.Sprintf("%%%s%%", function)
		}
		sb.WriteString(fmt.Sprintf("%s %s '%s'", CallerFunction.String(), operator.op(), function))
	})
}

func FilterByMessage(message string, operator StringOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		if operator == CONTAINS || operator == NOT_CONTAINS {
			message = fmt.Sprintf("%%%s%%", message)
		}
		sb.WriteString(fmt.Sprintf("%s %s '%s'", Message.String(), operator.op(), message))
	})
}

func FilterByTimestamp(timestamp string, operator QueryOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		if operator == CONTAINS || operator == NOT_CONTAINS {
			timestamp = fmt.Sprintf("%%%s%%", timestamp)
		}
		sb.WriteString(fmt.Sprintf("%s %s '%s'", Timestamp.String(), operator.op(), timestamp))
	})
}

func SortByLevel(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", Level.String(), order.op()))
	})
}

func SortByContext(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", Tags.String(), order.op()))
	})
}

func SortByCallerFile(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", CallerFile.String(), order.op()))
	})
}

func SortByCallerLine(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", Callerline.String(), order.op()))
	})
}

func SortByCallerFunction(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", CallerFunction.String(), order.op()))
	})
}

func SortByMessage(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", Message.String(), order.op()))
	})
}

func SortByTimestamp(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", Timestamp.String(), order.op()))
	})
}
