package postgres

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func normalizedServiceName(serviceName *string) (string, bool) {
	if serviceName == nil {
		return "", false
	}
	v := strings.TrimSpace(*serviceName)
	return v, v != ""
}

func addFilter(args []any, clauses []string, clause string, value any) ([]any, []string) {
	args = append(args, value)
	clauses = append(clauses, fmt.Sprintf(clause, len(args)))
	return args, clauses
}

func buildBaseFilters(userID *uuid.UUID, serviceName *string) ([]any, []string) {
	args := make([]any, 0, 2)
	clauses := make([]string, 0, 2)

	if userID != nil {
		args, clauses = addFilter(args, clauses, "user_id=$%d", *userID)
	}
	if serviceName != nil {
		if v, ok := normalizedServiceName(serviceName); ok {
			args, clauses = addFilter(args, clauses, "service_name ILIKE $%d", v)
		}
	}
	return args, clauses
}

func BuildCommonFiltersForTest(userID *uuid.UUID, serviceName *string) (string, []any) {
	args, clauses := buildBaseFilters(userID, serviceName)
	if len(clauses) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(clauses, " AND "), args
}
