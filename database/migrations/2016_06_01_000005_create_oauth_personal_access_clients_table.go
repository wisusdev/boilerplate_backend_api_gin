package migrations

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/database/generate_migrations"
	"boilerplate_backend_api_gin/core/database/schema"
)

type CreateOauthPersonalAccessClientsTable struct {
	generate_migrations.BaseMigration
}

func NewCreateOauthPersonalAccessClientsTable() *CreateOauthPersonalAccessClientsTable {
	return &CreateOauthPersonalAccessClientsTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_oauth_personal_access_clients_table",
			Timestamp: "2016_06_01_000005",
		},
	}
}

func (m *CreateOauthPersonalAccessClientsTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("oauth_personal_access_clients", func(table *schema.Blueprint) {
		table.Increments("id")
		table.BigInteger("client_id").Unsigned()
		table.Timestamps()
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreateOauthPersonalAccessClientsTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS oauth_personal_access_clients")
	return err
}
