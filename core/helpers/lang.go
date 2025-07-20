package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"semita/config"
	"strings"
)

var translations map[string]map[string]string
var defaultLang = "en"

func LoadTranslations() {
	translations = make(map[string]map[string]string)

	var langDir = "lang"

	var files []string
	var entries, err = os.ReadDir(langDir)
	if err != nil {
		Logs("ERROR", fmt.Sprintf("Failed to read lang directory: %v", err))
		return
	} else {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
				files = append(files, langDir+"/"+entry.Name())
			}
		}
	}

	for _, file := range files {
		lang := strings.TrimSuffix(strings.TrimPrefix(file, langDir+"/"), ".json")

		data, err := os.ReadFile(file)

		if err != nil {
			Logs("ERROR", fmt.Sprintf("Failed to read translation file: %s (%v)", file, err))
		}

		var m map[string]string

		if err := json.Unmarshal(data, &m); err != nil {
			Logs("ERROR", fmt.Sprintf("Failed to parse translation file: %s (%v)", file, err))
		}

		translations[lang] = m
	}

	fmt.Println("✅ Traducciones cargadas correctamente")
}

func Translate(key, lang string) string {
	lang = config.AppConfig().Lang
	if lang == "" {
		lang = defaultLang
	}

	if val, ok := translations[lang][key]; ok {
		return val
	}

	if val, ok := translations[defaultLang][key]; ok {
		return val
	}

	return key
}
