package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreatePasswordResetsTable struct {
	generate_migrations.BaseMigration
}

func NewCreatePasswordResetsTable() *CreatePasswordResetsTable {
	return &CreatePasswordResetsTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_password_resets_table",
			Timestamp: "2025_07_07_000001",
		},
	}
}

func (m *CreatePasswordResetsTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sqlQuery := schemaBuilder.Create("password_resets", func(table *schema.Blueprint) {
		table.String("email", 255)
		table.String("token", 255)
		table.DateTime("created_at")

		// Clave primaria compuesta
		table.Primary([]string{"email", "token"})
	})

	_, err := db.Exec(sqlQuery)
	return err
}

func (m *CreatePasswordResetsTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS password_resets")
	return err
}
