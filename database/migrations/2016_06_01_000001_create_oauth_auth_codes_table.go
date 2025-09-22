package migrations

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/database/generate_migrations"
	"boilerplate_backend_api_gin/core/database/schema"
)

type CreateOauthAuthCodesTable struct {
	generate_migrations.BaseMigration
}

func NewCreateOauthAuthCodesTable() *CreateOauthAuthCodesTable {
	return &CreateOauthAuthCodesTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_oauth_auth_codes_table",
			Timestamp: "2016_06_01_000001",
		},
	}
}

func (m *CreateOauthAuthCodesTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sql := schemaBuilder.Create("oauth_auth_codes", func(table *schema.Blueprint) {
		table.String("id", 100)
		table.Uuid("user_id").Index()
		table.BigInteger("client_id").Unsigned()
		table.Text("scopes").Nullable()
		table.Boolean("revoked")
		table.DateTime("expires_at").Nullable()
		table.Primary([]string{"id"})
	})

	_, err := db.Exec(sql)
	return err
}

func (m *CreateOauthAuthCodesTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS oauth_auth_codes")
	return err
}
