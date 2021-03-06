// Code generated by entc, DO NOT EDIT.

package todo

import (
	"time"

	"entgo.io/ent"
)

const (
	// Label holds the string label denoting the todo type in the database.
	Label = "todo"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldText holds the string denoting the text field in the database.
	FieldText = "text"
	// EdgeOwner holds the string denoting the owner edge name in mutations.
	EdgeOwner = "owner"
	// Table holds the table name of the todo in the database.
	Table = "todos"
	// OwnerTable is the table the holds the owner relation/edge.
	OwnerTable = "todos"
	// OwnerInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	OwnerInverseTable = "users"
	// OwnerColumn is the table column denoting the owner relation/edge.
	OwnerColumn = "user_todos"
)

// Columns holds all SQL columns for todo fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldText,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "todos"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_todos",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/roffe/gin-ent/ent/runtime"
//
var (
	Hooks  [1]ent.Hook
	Policy ent.Policy
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// TextValidator is a validator for the "text" field. It is called by the builders before save.
	TextValidator func(string) error
)
