package migrations

import (
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
	"semita/core/database/schema"
)

type CreateOAuthTokensTable struct {
	generate_migrations.BaseMigration
}

func NewCreateOAuthTokensTable() *CreateOAuthTokensTable {
	return &CreateOAuthTokensTable{
		BaseMigration: generate_migrations.BaseMigration{
			Name:      "create_oauth_tokens_table",
			Timestamp: "2025_07_06_000002",
		},
	}
}

func (m *CreateOAuthTokensTable) Up(db database_connections.SQLAdapter) error {
	// Usar Schema Builder para definir la tabla
	schemaBuilder := schema.NewSchema()

	sqlQuery := schemaBuilder.Create("oauth_tokens", func(table *schema.Blueprint) {
		table.Increments("id")
		table.UnsignedInteger("user_id").Nullable().Index()
		table.UnsignedInteger("client_id").Index()
		table.String("access_token", 512).Unique()
		table.String("refresh_token", 512).Unique()
		table.String("scopes", 255).Nullable()
		table.Boolean("revoked").Default(false)
		table.DateTime("expires_at")
		table.Timestamp("created_at").UseCurrent()
		table.Timestamp("updated_at").UseCurrent().OnUpdateCurrent()

		// Claves foráneas
		table.Foreign("user_id").References("id").On("users").OnDelete("CASCADE")
		table.Foreign("client_id").References("id").On("oauth_clients").OnDelete("CASCADE")
	})

	// Verificar el tipo del adaptador
	_, err := db.Exec(sqlQuery)
	return err
}

func (m *CreateOAuthTokensTable) Down(db database_connections.SQLAdapter) error {
	_, err := db.Exec("DROP TABLE IF EXISTS oauth_tokens")
	return err
}
