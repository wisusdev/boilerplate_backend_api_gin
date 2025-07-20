package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreateUserRolesTable struct {
	generate_migrations.BaseMigration
}

func NewCreateUserRolesTable() *CreateUserRolesTable {
	return &CreateUserRolesTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_user_roles_table",
			Timestamp: "2025_07_11_000003",
		},
	}
}

func (m *CreateUserRolesTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()
	sqlQuery := schemaBuilder.Create("user_roles", func(table *schema.Blueprint) {
		table.Increments("id")
		table.UnsignedInteger("user_id")
		table.UnsignedInteger("role_id")
		table.Timestamp("created_at").UseCurrent()
		table.Timestamp("updated_at").UseCurrent().OnUpdateCurrent()

		// Claves foráneas
		table.Foreign("user_id").References("id").On("users").OnDelete("CASCADE")
		table.Foreign("role_id").References("id").On("roles").OnDelete("CASCADE")

		// Índices únicos y regulares
		table.Unique([]string{"user_id", "role_id"})
		table.Index("user_id")
		table.Index("role_id")
	})

	_, err := db.Exec(sqlQuery)
	return err
}

func (m *CreateUserRolesTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS user_roles")
	return err
}
