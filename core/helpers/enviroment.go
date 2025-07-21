package helpers

import (
	"fmt"
	"os"
	"strings"
)

func UpdateEnvFile(key, value string) {
	file := ".env"
	input, err := os.ReadFile(file)

	if err != nil {
		fmt.Printf("No se pudo leer %s, solo se mostrará la clave.\n", file)
		return
	}

	lines := strings.Split(string(input), "\n")
	found := false

	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + value
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, key+"="+value)
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(file, []byte(output), 0644)

	if err != nil {
		fmt.Printf("❌ Error al actualizar %s: %v\n", file, err)
		return
	}

	fmt.Printf("✅ Clave: %s generada correctamente para %s en el archivo %s\n", value, key, file)
}
