package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreateUserPermissionsTable struct {
	generate_migrations.BaseMigration
}

func NewCreateUserPermissionsTable() *CreateUserPermissionsTable {
	return &CreateUserPermissionsTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_user_permissions_table",
			Timestamp: "2025_07_11_000005",
		},
	}
}

func (m *CreateUserPermissionsTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()
	sqlQuery := schemaBuilder.Create("user_permissions", func(table *schema.Blueprint) {
		table.Increments("id")
		table.UnsignedInteger("user_id")
		table.UnsignedInteger("permission_id")
		table.Timestamp("created_at").UseCurrent()
		table.Timestamp("updated_at").UseCurrent().OnUpdateCurrent()

		// Claves foráneas
		table.Foreign("user_id").References("id").On("users").OnDelete("CASCADE")
		table.Foreign("permission_id").References("id").On("permissions").OnDelete("CASCADE")

		// Índices únicos y regulares
		table.Unique([]string{"user_id", "permission_id"})
		table.Index("user_id")
		table.Index("permission_id")
	})

	_, err := db.Exec(sqlQuery)
	return err
}

func (m *CreateUserPermissionsTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS user_permissions")
	return err
}
