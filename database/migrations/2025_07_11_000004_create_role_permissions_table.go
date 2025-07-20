package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreateRolePermissionsTable struct {
	generate_migrations.BaseMigration
}

func NewCreateRolePermissionsTable() *CreateRolePermissionsTable {
	return &CreateRolePermissionsTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_role_permissions_table",
			Timestamp: "2025_07_11_000004",
		},
	}
}

func (m *CreateRolePermissionsTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()
	sqlQuery := schemaBuilder.Create("role_permissions", func(table *schema.Blueprint) {
		table.Increments("id")
		table.UnsignedInteger("role_id")
		table.UnsignedInteger("permission_id")
		table.Timestamp("created_at").UseCurrent()
		table.Timestamp("updated_at").UseCurrent().OnUpdateCurrent()

		// Claves foráneas
		table.Foreign("role_id").References("id").On("roles").OnDelete("CASCADE")
		table.Foreign("permission_id").References("id").On("permissions").OnDelete("CASCADE")

		// Índices únicos y regulares
		table.Unique([]string{"role_id", "permission_id"})
		table.Index("role_id")
		table.Index("permission_id")
	})

	_, err := db.Exec(sqlQuery)
	return err
}

func (m *CreateRolePermissionsTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS role_permissions")
	return err
}
