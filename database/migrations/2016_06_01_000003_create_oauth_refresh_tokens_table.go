package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreateOauthRefreshTokensTable struct {
	generate_migrations.BaseMigration
}

func NewCreateOauthRefreshTokensTable() *CreateOauthRefreshTokensTable {
	return &CreateOauthRefreshTokensTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_oauth_refresh_tokens_table",
			Timestamp: "2016_06_01_000003",
		},
	}
}

func (m *CreateOauthRefreshTokensTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("oauth_refresh_tokens", func(table *schema.Blueprint) {
		table.String("id", 100)
		table.String("access_token_id", 100).Index()
		table.Boolean("revoked")
		table.DateTime("expires_at").Nullable()
		table.Primary([]string{"id"})
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreateOauthRefreshTokensTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS oauth_refresh_tokens")
	return err
}
