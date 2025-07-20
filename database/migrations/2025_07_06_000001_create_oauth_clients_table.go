package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreateOAuthClientsTable struct {
	generate_migrations.BaseMigration
}

func NewCreateOAuthClientsTable() *CreateOAuthClientsTable {
	return &CreateOAuthClientsTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_oauth_clients_table",
			Timestamp: "2025_07_06_000001",
		},
	}
}

func (m *CreateOAuthClientsTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("oauth_clients", func(table *schema.Blueprint) {
		table.Increments("id")
		table.String("name", 255)
		table.String("client_id", 100).Unique()
		table.String("client_secret", 255)
		table.String("redirect_uri", 255).Nullable()
		table.String("grant_types", 255).Nullable()
		table.String("scopes", 255).Nullable()
		table.Timestamp("created_at").UseCurrent()
		table.Timestamp("updated_at").UseCurrent().OnUpdateCurrent()
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreateOAuthClientsTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS oauth_clients")
	return err
}
