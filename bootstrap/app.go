package bootstrap

import (
	"fmt"
	"os"
	"semita/core/cli"
	"semita/core/commands"
	"semita/core/helpers"
)

func Commands() {
	cli.RegisterCommand(commands.MigrateCmd)
	cli.RegisterCommand(commands.MigrateFreshCmd)
	cli.RegisterCommand(commands.MigrateRollbackCmd)
	cli.RegisterCommand(commands.MakeMigrationCmd)
	cli.RegisterCommand(commands.KeyGenerateCmd)
	cli.RegisterCommand(commands.OauthKeysCmd)
	cli.RegisterCommand(commands.OauthClientCmd)
	cli.RegisterCommand(commands.SeedAllCommand)
	cli.RegisterCommand(commands.SeedRunCommand)
}

// Execute inicializa y ejecuta los comandos
func Execute(args []string) {
	Commands() // Registrar los comandos
	if args[0] == "help" {
		cli.ListCommands(nil)
		return
	}

	commandName := args[0]
	commandArgs := args[1:]

	cmd, ok := cli.Commands[commandName]
	if !ok {
		helpers.Logs("ERROR", "Comando no encontrado: "+commandName)
		fmt.Printf("Comando '%s' no encontrado. Ejecuta 'go run main.go list' para ver los comandos disponibles.\n", commandName)
		os.Exit(1)
	}

	cmd.Execute(commandArgs)
}
