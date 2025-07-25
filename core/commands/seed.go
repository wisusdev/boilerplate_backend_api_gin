package commands

import (
	"fmt"
	"log"
	"semita/core/helpers"
	"semita/database/seeders"

	"github.com/spf13/cobra"
)

// SeedAllCommand ejecuta todos los generate_seeders
var SeedAllCommand = &cobra.Command{
	Use:   "db:seed",
	Short: "Ejecuta todos los generate_seeders",
	Long:  "Execute all registered generate_seeders in the correct dependency order.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runAllSeeders(); err != nil {
			log.Fatalf("Error running all generate_seeders: %v", err)
		}
	},
}

// SeedRunCommand ejecuta un seeder específico
var SeedRunCommand = &cobra.Command{
	Use:   "run:seed [seeder_name]",
	Short: "Ejecuta un seeder específico",
	Long:  "Execute a specific seeder by name. Dependencies will be run first if needed.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
