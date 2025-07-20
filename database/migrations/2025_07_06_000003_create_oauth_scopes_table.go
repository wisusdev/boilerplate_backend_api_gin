package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreateOAuthScopesTable struct {
	generate_migrations.BaseMigration
}

func NewCreateOAuthScopesTable() *CreateOAuthScopesTable {
	return &CreateOAuthScopesTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_oauth_scopes_table",
			Timestamp: "2025_07_06_000003",
		},
	}
}

func (m *CreateOAuthScopesTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sqlQuery := schemaBuilder.Create("oauth_scopes", func(table *schema.Blueprint) {
		table.Increments("id")
		table.String("name", 100).Unique()
		table.String("description", 255).Nullable()
		table.Timestamp("created_at").UseCurrent()
		table.Timestamp("updated_at").UseCurrent().OnUpdateCurrent()
	})

	_, err := db.Exec(sqlQuery)
	return err
}

func (m *CreateOAuthScopesTable) Down(db database_connections.SQLAdapter) error {

	_, err := db.Exec("DROP TABLE IF EXISTS oauth_scopes")
	return err
}
