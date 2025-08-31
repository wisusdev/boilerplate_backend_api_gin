package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreateOauthAccessTokensTable struct {
	generate_migrations.BaseMigration
}

func NewCreateOauthAccessTokensTable() *CreateOauthAccessTokensTable {
	return &CreateOauthAccessTokensTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_oauth_access_tokens_table",
			Timestamp: "2016_06_01_000002",
		},
	}
}

func (m *CreateOauthAccessTokensTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("oauth_access_tokens", func(table *schema.Blueprint) {
		table.String("id", 100)
		table.Uuid("user_id").Nullable().Index()
		table.BigInteger("client_id").Unsigned()
		table.String("name", 255).Nullable()
		table.Text("scopes").Nullable()
		table.Boolean("revoked")
		table.Timestamps()
		table.DateTime("expires_at").Nullable()
		table.Primary([]string{"id"})
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreateOauthAccessTokensTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS oauth_access_tokens")
	return err
}
