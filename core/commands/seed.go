package commands

import (
	"fmt"
	"log"
	"semita/core/cli"
	"semita/core/helpers"
	"semita/database/seeders"
)

// SeedAllCommand ejecuta todos los generate_seeders
var SeedAllCommand = cli.Command{
	Name:        "db:seed",
	Description: "Ejecuta todos los generate_seeders",
	Execute: func(args []string) {
		if err := runAllSeeders(); err != nil {
			log.Fatalf("Error running all generate_seeders: %v", err)
		}
	},
}

// SeedRunCommand ejecuta un seeder específico
var SeedRunCommand = cli.Command{
	Name:        "run:seed [seeder_name]",
	Description: "Ejecuta un seeder específico",
	Execute: func(args []string) {
		runSpecificSeeder(args[0])
	},
}

// runAllSeeders ejecuta todos los generate_seeders
func runAllSeeders() error {
	manager := seeders.CreateSeederManager()
	err := manager.RunAllSeeders()
	if err != nil {
		log.Fatalf("Error running all generate_seeders: %v", err)
		return err
	}

	return nil
}

// runSpecificSeeder ejecuta un seeder específico
func runSpecificSeeder(seederName string) {
	manager := seeders.CreateSeederManager()
	err := manager.RunSeeder(seederName)
	if err != nil {
		helpers.Logs("ERROR", fmt.Sprintf("%v", err))
		log.Fatalf("Error running seeder '%s': %v", seederName, err)
	}

	log.Printf("=== Seeder '%s' Completed Successfully ===", seederName)
}
