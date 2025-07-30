package cli

import (
	"fmt"
	"semita/core/helpers"
	"sort"
)

// Command define la estructura de un comando de consola
type Command struct {
	Name        string
	Description string
	Execute     func(args []string)
}

// Commands es nuestro registro de todos los comandos de la aplicación
var Commands = make(map[string]Command)

// RegisterCommand añade un nuevo comando a nuestro registro.
// Es "thread-safe" por si se quisiera usar en un futuro con goroutines.
func RegisterCommand(cmd Command) {
	// En una aplicación real, aquí se usaría un `sync.Mutex` si se registran comandos concurrentemente.
	// Para un CLI simple, no es necesario.
	Commands[cmd.Name] = cmd
}

// Comando para listar todos los demás comandos
func ListCommands(args []string) {
	fmt.Println(helpers.ColorMagenta("Semita CLI (versión casera)"))
	fmt.Println("---------------------------")
	fmt.Println(helpers.ColorYellow("Uso: go run main.go <comando> [argumentos]"))
	fmt.Println("\nComandos disponibles:")

	// Para que la lista salga ordenada alfabéticamente
	keys := make([]string, 0, len(Commands))
	for k := range Commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		// Usamos fmt.Printf para alinear las descripciones y que se vea limpio
		fmt.Printf("  %s %s\n", helpers.ColorGreen(fmt.Sprintf("%-30s", name)), Commands[name].Description)
	}
}
