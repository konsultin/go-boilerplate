package query

import (
	"encoding/json"
	"fmt"

	"github.com/konsultin/project-goes-here/libs/sqlk/schema"
)

// AuditData holds audit field information
type AuditData struct {
	CreatedAt interface{} // timestamp or nil
	UpdatedAt interface{} // timestamp or nil
	CreatedBy interface{} // Subject struct/map to be marshaled to JSON
	UpdatedBy interface{} // Subject struct/map to be marshaled to JSON
}

// PrepareAuditFields prepares audit field values for INSERT/UPDATE
// Returns map of column name to value, with JSON marshaling for Subject fields
func PrepareAuditFields(s *schema.Schema, data *AuditData, isUpdate bool) (map[string]interface{}, error) {
	if !s.AuditFields() {
		return nil, nil
	}

	result := make(map[string]interface{})

	if !isUpdate {
		// INSERT: Set created_at and created_by
		if data.CreatedAt != nil {
			result[s.CreatedAtColumn()] = data.CreatedAt
		} else {
			// Auto-populate with NOW()
			result[s.CreatedAtColumn()] = "NOW()"
		}

		if data.CreatedBy != nil {
			// Marshal to JSON
			jsonBytes, err := json.Marshal(data.CreatedBy)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal created_by to JSON: %w", err)
			}
			result[s.CreatedByColumn()] = string(jsonBytes)
		}
	}

	// UPDATE: Set updated_at and updated_by
	if data.UpdatedAt != nil {
		result[s.UpdatedAtColumn()] = data.UpdatedAt
	} else {
		// Auto-populate with NOW()
		result[s.UpdatedAtColumn()] = "NOW()"
	}

	if data.UpdatedBy != nil {
		// Marshal to JSON
		jsonBytes, err := json.Marshal(data.UpdatedBy)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal updated_by to JSON: %w", err)
		}
		result[s.UpdatedByColumn()] = string(jsonBytes)
	}

	return result, nil
}

// AddAuditFieldsToInsert adds audit fields to INSERT column list
func AddAuditFieldsToInsert(s *schema.Schema, columns []string, data *AuditData) ([]string, error) {
	if !s.AuditFields() {
		return columns, nil
	}

	auditFields, err := PrepareAuditFields(s, data, false)
	if err != nil {
		return nil, err
	}

	for col := range auditFields {
		// Check if column is not already in the list
		found := false
		for _, existing := range columns {
			if existing == col {
				found = true
				break
			}
		}
		if !found {
			columns = append(columns, col)
		}
	}

	return columns, nil
}

// AddAuditFieldsToUpdate adds audit fields to UPDATE column list
func AddAuditFieldsToUpdate(s *schema.Schema, columns []string, data *AuditData) ([]string, error) {
	if !s.AuditFields() {
		return columns, nil
	}

	auditFields, err := PrepareAuditFields(s, data, true)
	if err != nil {
		return nil, err
	}

	for col := range auditFields {
		// Check if column is not already in the list
		found := false
		for _, existing := range columns {
			if existing == col {
				found = true
				break
			}
		}
		if !found {
			columns = append(columns, col)
		}
	}

	return columns, nil
}

// Note: The actual integration with Insert and Update builders
// should be done by the developer when calling these builders
// by including the audit columns and values in their arguments
