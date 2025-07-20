package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreateRolesTable struct {
	generate_migrations.BaseMigration
}

func NewCreateRolesTable() *CreateRolesTable {
	return &CreateRolesTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_roles_table",
			Timestamp: "2025_07_11_000001",
		},
	}
}

func (m *CreateRolesTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sqlQuery := schemaBuilder.Create("roles", func(table *schema.Blueprint) {
		table.Increments("id")
		table.String("name", 255).Unique()
		table.String("guard_name", 255).Default("web")
		table.Text("description").Nullable()
		table.Timestamp("created_at").UseCurrent()
		table.Timestamp("updated_at").UseCurrent().OnUpdateCurrent()

		// Índices
		table.Index("name")
		table.Index("guard_name")
	})

	_, err := db.Exec(sqlQuery)
	return err
}

func (m *CreateRolesTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS roles")
	return err
}
