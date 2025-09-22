package migrations

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/database/generate_migrations"
	"boilerplate_backend_api_gin/core/database/schema"
)

type CreateOAuthClientsTable struct {
	generate_migrations.BaseMigration
}

func NewCreateOAuthClientsTable() *CreateOAuthClientsTable {
	return &CreateOAuthClientsTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_oauth_clients_table",
			Timestamp: "2016_06_01_000004",
		},
	}
}

func (m *CreateOAuthClientsTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("oauth_clients", func(table *schema.Blueprint) {
		table.Increments("id")
		table.BigInteger("user_id").Unsigned().Nullable().Index()
		table.String("name", 255)
		table.String("secret", 100).Nullable()
		table.String("provider", 255).Nullable()
		table.Text("redirect")
		table.Boolean("personal_access_client")
		table.Boolean("password_client")
		table.Boolean("revoked")
		table.Timestamps()
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreateOAuthClientsTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS oauth_clients")
	return err
}
