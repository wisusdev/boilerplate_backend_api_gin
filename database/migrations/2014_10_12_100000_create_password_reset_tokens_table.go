package migrations

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/database/generate_migrations"
	"boilerplate_backend_api_gin/core/database/schema"
)

type CreatePasswordResetTokensTable struct {
	generate_migrations.BaseMigration
}

func NewCreatePasswordResetTokensTable() *CreatePasswordResetTokensTable {
	return &CreatePasswordResetTokensTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_password_reset_tokens_table",
			Timestamp: "2014_10_12_100000",
		},
	}
}

func (m *CreatePasswordResetTokensTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("password_reset_tokens", func(table *schema.Blueprint) {
		table.String("email", 255)
		table.String("token", 255)
		table.Timestamp("created_at").Nullable()
		table.Primary([]string{"email"})
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreatePasswordResetTokensTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS password_reset_tokens")
	return err
}
