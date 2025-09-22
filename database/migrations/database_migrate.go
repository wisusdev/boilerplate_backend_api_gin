package migrations

import (
	"boilerplate_backend_api_gin/core/database/database_connections"
	"boilerplate_backend_api_gin/core/database/generate_migrations"
	"fmt"
)

func WithMigrator(action func(migrator *generate_migrations.Migrator)) {
	fmt.Println("🔌 Conectando a la base de datos...")
	db := database_connections.DatabaseConnectSQL()
	defer func(db database_connections.SQLAdapter) {
		err := db.Close()
		if err != nil {
			fmt.Println("❌ Error al cerrar la conexión a la base de datos:", err)
		} else {
			fmt.Println("✅ Conexión a la base de datos cerrada correctamente.")
		}
	}(db)

	migrator := generate_migrations.NewMigrator(db)

	// Registrar migraciones en el orden cronológico de Laravel
	migrator.Register(NewCreateUsersTable())
	migrator.Register(NewCreatePasswordResetTokensTable())
	migrator.Register(NewCreateOauthAuthCodesTable())
	migrator.Register(NewCreateOauthAccessTokensTable())
	migrator.Register(NewCreateOauthRefreshTokensTable())
	migrator.Register(NewCreateOAuthClientsTable())
	migrator.Register(NewCreateOauthPersonalAccessClientsTable())
	migrator.Register(NewCreateFailedJobsTable())
	migrator.Register(NewCreatePersonalAccessTokensTable())
	migrator.Register(NewCreatePermissionTablesTable())

	fmt.Println("🚀 Ejecutando acción del migrator...")
	action(migrator)
}
