package migrations

import (
	"fmt"
	"semita/core/database/database_connections"
	"semita/core/database/generate_migrations"
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

	migrator.Register(NewCreateUsersTable())
	migrator.Register(NewCreatePasswordResetsTable())
	migrator.Register(NewCreateOAuthClientsTable())
	migrator.Register(NewCreateOAuthTokensTable())
	migrator.Register(NewCreateOAuthScopesTable())
	migrator.Register(NewCreateRolesTable())
	migrator.Register(NewCreatePermissionsTable())
	migrator.Register(NewCreateUserRolesTable())
	migrator.Register(NewCreateRolePermissionsTable())
	migrator.Register(NewCreateUserPermissionsTable())

	fmt.Println("🚀 Ejecutando acción del migrator...")
	action(migrator)
}
