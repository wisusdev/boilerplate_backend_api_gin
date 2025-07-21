package bootstrap

import (
	"semita/core/commands"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "semita",
	Short: "Semita CLI",
	// Si no hay subcomando, muestra la ayuda
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Commands() {
	RootCmd.AddCommand(commands.MigrateCmd)
	RootCmd.AddCommand(commands.MigrateFreshCmd)
	RootCmd.AddCommand(commands.MigrateRollbackCmd)
	RootCmd.AddCommand(commands.MakeMigrationCmd)
	RootCmd.AddCommand(commands.KeyGenerateCmd)
	RootCmd.AddCommand(commands.OauthKeysCmd)
	RootCmd.AddCommand(commands.OauthClientCmd)
	RootCmd.AddCommand(commands.SeedAllCommand)
	RootCmd.AddCommand(commands.SeedRunCommand)
}

// Execute inicializa y ejecuta los comandos
func Execute() {
	Commands() // Registrar los comandos
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}
