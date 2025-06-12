package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// DISCLAIMER THIS FILE IS COMPLETELY VIBE CODED YOU HAVE BEEN WARNED
// DISCLAIMER THIS FILE IS COMPLETELY VIBE CODED YOU HAVE BEEN WARNED
// DISCLAIMER THIS FILE IS COMPLETELY VIBE CODED YOU HAVE BEEN WARNED
// DISCLAIMER THIS FILE IS COMPLETELY VIBE CODED YOU HAVE BEEN WARNED
// DISCLAIMER THIS FILE IS COMPLETELY VIBE CODED YOU HAVE BEEN WARNED

// MockDatabase is a mock implementation of IDatabase for testing
type MockDatabase struct {
	ctx           context.Context
	errorToReturn error
	queryParams   map[string]any
	records       map[reflect.Type][]any // Type -> Slice of records
	limit         int
	offset        int
	count         int64
	inTransaction bool
	committed     bool
	rolledBack    bool
	nextID        int64 // For auto-incrementing IDs
}

// NewMockDatabase creates a new instance of MockDatabase
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		queryParams: make(map[string]any),
		records:     make(map[reflect.Type][]any),
		nextID:      1,
	}
}

func (m *MockDatabase) Where(query any, args ...any) IDatabase {
	m.queryParams["where_query"] = query
	if len(args) > 0 {
		m.queryParams["where_args"] = args
	}
	return m
}

func (m *MockDatabase) First(dest any, conds ...any) IDatabase {
	if m.errorToReturn != nil {
		return m
	}

	if len(conds) > 0 {
		m.queryParams["first_conds"] = conds
	}

	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		m.errorToReturn = errors.New("destination must be a pointer")
		return m
	}

	destType := destValue.Elem().Type()
	records, exists := m.records[destType]
	if !exists || len(records) == 0 {
		m.errorToReturn = errors.New("record not found")
		return m
	}

	// Apply where conditions if they exist
	filtered := m.applyConditions(records)
	if len(filtered) == 0 {
		m.errorToReturn = errors.New("record not found")
		return m
	}

	// Get first record
	srcValue := reflect.ValueOf(filtered[0])
	if destValue.Elem().Type() == srcValue.Type() {
		destValue.Elem().Set(srcValue)
	}

	return m
}

func (m *MockDatabase) Create(value any) IDatabase {
	if m.errorToReturn != nil {
		return m
	}

	valueType := reflect.TypeOf(value)
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}

	// Check if value has an ID field and set it if zero
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		idField := val.FieldByName("ID")
		if idField.IsValid() && idField.CanSet() {
			switch idField.Kind() {
			case reflect.Int, reflect.Int64:
				if idField.Int() == 0 {
					idField.SetInt(m.nextID)
					m.nextID++
				}
			case reflect.Uint, reflect.Uint64:
				if idField.Uint() == 0 {
					idField.SetUint(uint64(m.nextID))
					m.nextID++
				}
			}
		}
	}

	m.records[valueType] = append(m.records[valueType], value)
	m.queryParams["create_value"] = value
	return m
}

func (m *MockDatabase) Delete(value any, conds ...any) IDatabase {
	if m.errorToReturn != nil {
		return m
	}

	if len(conds) > 0 {
		m.queryParams["delete_conds"] = conds
	}

	valueType := reflect.TypeOf(value)
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}

	records, exists := m.records[valueType]
	if !exists {
		m.errorToReturn = errors.New("no records of this type")
		return m
	}

	// Check if we have any where conditions
	hasWhere := false
	if whereQuery, exists := m.queryParams["where_query"]; exists {
		switch v := whereQuery.(type) {
		case string:
			if v != "" {
				hasWhere = true
			}
		case map[string]any:
			if len(v) > 0 {
				hasWhere = true
			}
		}
	}

	hasArgs := false
	if whereArgs, exists := m.queryParams["where_args"]; exists {
		switch v := whereArgs.(type) {
		case []any:
			if len(v) > 0 {
				hasArgs = true
			}
		}
	}

	// If no conditions, delete all records of this type
	if !hasWhere && !hasArgs {
		delete(m.records, valueType)
		return m
	}

	// Apply conditions to find records to delete
	filtered := m.applyConditions(records)

	// Create a map of records to delete for quick lookup
	toDelete := make(map[any]bool)
	for _, rec := range filtered {
		toDelete[rec] = true
	}

	// Filter out records to be deleted
	var remaining []any
	for _, record := range records {
		if !toDelete[record] {
			remaining = append(remaining, record)
		}
	}

	m.records[valueType] = remaining
	return m
}

func (m *MockDatabase) Model(value any) IDatabase {
	m.queryParams["model_value"] = value
	return m
}

func (m *MockDatabase) Updates(values any) IDatabase {
	if m.errorToReturn != nil {
		return m
	}

	m.queryParams["updates_values"] = values

	// Get the model type from Model() call
	modelValue, ok := m.queryParams["model_value"]
	if !ok {
		m.errorToReturn = errors.New("model not specified")
		return m
	}

	modelType := reflect.TypeOf(modelValue)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	records, exists := m.records[modelType]
	if !exists {
		m.errorToReturn = errors.New("no records of this type")
		return m
	}

	// Apply updates to all matching records
	filtered := m.applyConditions(records)
	updatesMap := m.structToMap(values)

	for i, record := range m.records[modelType] {
		for _, filteredRec := range filtered {
			if reflect.DeepEqual(record, filteredRec) {
				// Apply updates
				recValue := reflect.ValueOf(record)
				if recValue.Kind() == reflect.Ptr {
					recValue = recValue.Elem()
				}

				for k, v := range updatesMap {
					field := recValue.FieldByName(k)
					if field.IsValid() && field.CanSet() {
						field.Set(reflect.ValueOf(v))
					}
				}

				m.records[modelType][i] = record
				break
			}
		}
	}

	return m
}

func (m *MockDatabase) Limit(limit int) IDatabase {
	m.limit = limit
	return m
}

func (m *MockDatabase) Offset(offset int) IDatabase {
	m.offset = offset
	return m
}

func (m *MockDatabase) Order(value any) IDatabase {
	m.queryParams["order_value"] = value
	return m
}

func (m *MockDatabase) Count(count *int64) IDatabase {
	if m.errorToReturn != nil {
		return m
	}

	modelValue, ok := m.queryParams["model_value"]
	if !ok {
		m.errorToReturn = errors.New("model not specified")
		return m
	}

	modelType := reflect.TypeOf(modelValue)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	records, exists := m.records[modelType]
	if !exists {
		*count = 0
		return m
	}

	filtered := m.applyConditions(records)
	*count = int64(len(filtered))
	return m
}

func (m *MockDatabase) WithContext(ctx context.Context) IDatabase {
	m.ctx = ctx
	return m
}

func (m *MockDatabase) Begin(opts ...*sql.TxOptions) IDatabase {
	m.inTransaction = true
	m.committed = false
	m.rolledBack = false
	return m
}

func (m *MockDatabase) Commit() IDatabase {
	if m.inTransaction {
		m.committed = true
		m.inTransaction = false
	} else {
		m.errorToReturn = errors.New("no transaction in progress")
	}
	return m
}

func (m *MockDatabase) Rollback() IDatabase {
	if m.inTransaction {
		m.rolledBack = true
		m.inTransaction = false
	} else {
		m.errorToReturn = errors.New("no transaction in progress")
	}
	return m
}

func (m *MockDatabase) DB() (*sql.DB, error) {
	return nil, fmt.Errorf("not implemented in mock")
}

func (m *MockDatabase) AutoMigrate(dst ...any) error {
	if m.errorToReturn != nil {
		return m.errorToReturn
	}
	m.queryParams["auto_migrate_dst"] = dst

	// Initialize empty slices for each type
	for _, d := range dst {
		t := reflect.TypeOf(d)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if _, exists := m.records[t]; !exists {
			m.records[t] = []any{}
		}
	}

	return nil
}

func (m *MockDatabase) Error() error {
	err := m.errorToReturn
	m.errorToReturn = nil // Reset error after reading
	return err
}

// Helper methods

func (m *MockDatabase) applyConditions(records []any) []any {
	var filtered []any

	whereQuery, hasWhere := m.queryParams["where_query"]
	whereArgs, hasArgs := m.queryParams["where_args"]

	if !hasWhere && !hasArgs {
		return records
	}

	for _, record := range records {
		match := true

		// Simple equality matching for demonstration
		// In a real implementation, you'd parse the whereQuery
		if hasWhere {
			switch q := whereQuery.(type) {
			case string:
				// Very basic string matching for example purposes
				if strings.Contains(q, "=") {
					parts := strings.Split(q, "=")
					if len(parts) == 2 {
						fieldName := strings.TrimSpace(parts[0])
						expectedValue := strings.Trim(parts[1], "' \"")

						recValue := reflect.ValueOf(record)
						if recValue.Kind() == reflect.Ptr {
							recValue = recValue.Elem()
						}

						field := recValue.FieldByName(fieldName)
						if field.IsValid() {
							actualValue := fmt.Sprintf("%v", field.Interface())
							if actualValue != expectedValue {
								match = false
							}
						}
					}
				}
			case map[string]any:
				// Handle map conditions
				recValue := reflect.ValueOf(record)
				if recValue.Kind() == reflect.Ptr {
					recValue = recValue.Elem()
				}

				for k, v := range q {
					field := recValue.FieldByName(k)
					if !field.IsValid() || !reflect.DeepEqual(field.Interface(), v) {
						match = false
						break
					}
				}
			}
		}

		if match && hasArgs {
			// Handle where args (very simplified)
			recValue := reflect.ValueOf(record)
			if recValue.Kind() == reflect.Ptr {
				recValue = recValue.Elem()
			}

			args := whereArgs.([]any)
			for i := 0; i < len(args); i += 2 {
				fieldName, ok := args[i].(string)
				if !ok {
					match = false
					break
				}

				expectedValue := args[i+1]
				field := recValue.FieldByName(fieldName)
				if !field.IsValid() || !reflect.DeepEqual(field.Interface(), expectedValue) {
					match = false
					break
				}
			}
		}

		if match {
			filtered = append(filtered, record)
		}
	}

	// Apply limit and offset
	if m.offset > 0 && len(filtered) > m.offset {
		filtered = filtered[m.offset:]
	}
	if m.limit > 0 && len(filtered) > m.limit {
		filtered = filtered[:m.limit]
	}

	return filtered
}

func (m *MockDatabase) structToMap(s any) map[string]any {
	result := make(map[string]any)
	val := reflect.ValueOf(s)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return result
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}
		result[field.Name] = fieldValue.Interface()
	}

	return result
}

func (m *MockDatabase) Find(dest any, conds ...any) IDatabase {
	if m.errorToReturn != nil {
		return m
	}

	// Store conditions if provided
	if len(conds) > 0 {
		m.queryParams["find_conds"] = conds
	}

	// Get the type of the destination slice elements
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		m.errorToReturn = errors.New("destination must be a pointer to a slice")
		return m
	}

	sliceElemType := destValue.Elem().Type().Elem()
	if sliceElemType.Kind() == reflect.Ptr {
		sliceElemType = sliceElemType.Elem()
	}

	// Get records of the matching type
	records, exists := m.records[sliceElemType]
	if !exists {
		// Return empty slice if no records exist
		destValue.Elem().Set(reflect.MakeSlice(destValue.Elem().Type(), 0, 0))
		return m
	}

	// Apply conditions if they exist
	filtered := m.applyConditions(records)

	// Create a new slice of the appropriate type
	resultSlice := reflect.MakeSlice(destValue.Elem().Type(), len(filtered), len(filtered))

	// Copy matching records to the destination slice
	for i, record := range filtered {
		recordValue := reflect.ValueOf(record)
		if destValue.Elem().Type().Elem().Kind() == reflect.Ptr {
			// If destination is pointer slice, create new pointer
			newPtr := reflect.New(sliceElemType)
			newPtr.Elem().Set(recordValue)
			resultSlice.Index(i).Set(newPtr)
		} else {
			// If destination is value slice, set value directly
			resultSlice.Index(i).Set(recordValue)
		}
	}

	destValue.Elem().Set(resultSlice)
	return m
}

func (m *MockDatabase) Select(query any, args ...any) IDatabase {
	if m.errorToReturn != nil {
		return m
	}

	// Store the select fields/query
	m.queryParams["select_query"] = query
	if len(args) > 0 {
		m.queryParams["select_args"] = args
	}

	// In a real mock, you might want to implement actual field selection filtering
	// For this basic mock, we'll just record the select operation
	return m
}

func (m *MockDatabase) Exec(query any, args ...any) IDatabase {
	if m.errorToReturn != nil {
		return m
	}

	// Store the exec query and args
	m.queryParams["exec_query"] = query
	if len(args) > 0 {
		m.queryParams["exec_args"] = args
	}

	// Handle different query types
	switch q := query.(type) {
	case string:
		// Simple SQL command parsing for testing purposes
		if strings.HasPrefix(strings.ToUpper(q), "DELETE FROM") {
			// Handle DELETE queries
			tableName := strings.TrimSpace(strings.TrimPrefix(strings.ToUpper(q), "DELETE FROM"))
			tableName = strings.Fields(tableName)[0] // Get first word (table name)

			// Find the type by table name (simplified for mock)
			for t := range m.records {
				if strings.EqualFold(t.Name(), tableName) {
					// Delete all records of this type
					delete(m.records, t)
					break
				}
			}
		}
		// Could add handling for other SQL commands here
	}

	return m
}

// Test Helpers

// SetError sets an error to be returned by subsequent operations
func (m *MockDatabase) SetError(err error) {
	m.errorToReturn = err
}

// SetCount sets the count value to be returned by Count()
func (m *MockDatabase) SetCount(count int64) {
	m.count = count
}

// AddRecord adds a record to the mock database
func (m *MockDatabase) AddRecord(record any) {
	recordType := reflect.TypeOf(record)
	if recordType.Kind() == reflect.Ptr {
		recordType = recordType.Elem()
	}
	m.records[recordType] = append(m.records[recordType], record)
}

// ClearRecords clears all records in the mock database
func (m *MockDatabase) ClearRecords() {
	m.records = make(map[reflect.Type][]any)
	m.nextID = 1
}

// GetRecords returns all records of a specific type
func (m *MockDatabase) GetRecords(model any) []any {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	return m.records[modelType]
}

// GetQueryParam returns a query parameter that was set
func (m *MockDatabase) GetQueryParam(key string) any {
	return m.queryParams[key]
}

// IsCommitted returns whether Commit() was called
func (m *MockDatabase) IsCommitted() bool {
	return m.committed
}

// IsRolledBack returns whether Rollback() was called
func (m *MockDatabase) IsRolledBack() bool {
	return m.rolledBack
}

// IsInTransaction returns whether a transaction is in progress
func (m *MockDatabase) IsInTransaction() bool {
	return m.inTransaction
}
