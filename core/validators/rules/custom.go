package validators

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

// FileRule válida archivos subidos
type FileRule struct {
	MaxSize      int64    // Tamaño máximo en bytes
	AllowedMimes []string // Tipos MIME permitidos
	AllowedExits []string // Extensiones permitidas
}

func (r *FileRule) Name() string {
	return "file"
}

func (r *FileRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return fmt.Errorf("the field must be a file")
	}

	if err := r.validateFileSize(file); err != nil {
		return err
	}
	if err := r.validateFileExtension(file); err != nil {
		return err
	}
	if err := r.validateFileMime(file); err != nil {
		return err
	}

	return nil
}

func (r *FileRule) validateFileSize(file *multipart.FileHeader) error {
	if r.MaxSize > 0 && file.Size > r.MaxSize {
		return fmt.Errorf("the file may not be greater than %d bytes", r.MaxSize)
	}
	return nil
}

func (r *FileRule) validateFileExtension(file *multipart.FileHeader) error {
	if len(r.AllowedExits) == 0 {
		return nil
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	ext = strings.TrimPrefix(ext, ".")
	for _, allowedExt := range r.AllowedExits {
		if ext == strings.ToLower(allowedExt) {
			return nil
		}
	}
	return fmt.Errorf("the file must be of type: %s", strings.Join(r.AllowedExits, ", "))
}

func (r *FileRule) validateFileMime(file *multipart.FileHeader) error {
	if len(r.AllowedMimes) == 0 {
		return nil
	}
	f, err := file.Open()
	if err != nil {
		return fmt.Errorf("unable to read file")
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(f)

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return fmt.Errorf("unable to read file content")
	}
	contentType := file.Header.Get("Content-Type")
	for _, allowedMime := range r.AllowedMimes {
		if contentType == allowedMime {
			return nil
		}
	}
	return fmt.Errorf("the file must be of type: %s", strings.Join(r.AllowedMimes, ", "))
}

// ImageRule valida archivos de imagen con dimensiones
type ImageRule struct {
	FileRule
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

func (r *ImageRule) Name() string {
	return "image"
}

func (r *ImageRule) Validate(value interface{}, data map[string]interface{}) error {
	// Primero validar como archivo
	if err := r.FileRule.Validate(value, data); err != nil {
		return err
	}

	// Aquí podrías añadir validación de dimensiones usando una librería de imágenes
	// Por simplicidad, solo validamos que sea una imagen por extensión
	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return fmt.Errorf("the field must be an image file")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	imageExits := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}

	isImage := false
	for _, imgExt := range imageExits {
		if ext == imgExt {
			isImage = true
			break
		}
	}

	if !isImage {
		return fmt.Errorf("the field must be an image")
	}

	return nil
}

// ArrayRule válida arrays con reglas específicas para cada elemento
type ArrayRule struct {
	MinItems int
	MaxItems int
	// ItemRules []Rule - Comentado temporalmente, se puede implementar más adelante
}

func (r *ArrayRule) Name() string {
	return "array"
}

func (r *ArrayRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	arr, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("the field must be an array")
	}

	// Validar cantidad de elementos
	if r.MinItems > 0 && len(arr) < r.MinItems {
		return fmt.Errorf("the field must have at least %d items", r.MinItems)
	}

	if r.MaxItems > 0 && len(arr) > r.MaxItems {
		return fmt.Errorf("the field may not have more than %d items", r.MaxItems)
	}

	// Validar cada elemento del array
	// TODO: Implementar validación de elementos individuales cuando sea necesario
	// for i, item := range arr {
	//     for _, rule := range r.ItemRules {
	//         if err := rule.Validate(item, data); err != nil {
	//             return fmt.Errorf("item at index %d: %s", i, err.Error())
	//         }
	//     }
	// }

	return nil
}

// CustomRule permite definir reglas personalizadas
type CustomRule struct {
	RuleName     string
	ValidateFunc func(value interface{}, data map[string]interface{}) error
}

func (r *CustomRule) Name() string {
	return r.RuleName
}

func (r *CustomRule) Validate(value interface{}, data map[string]interface{}) error {
	return r.ValidateFunc(value, data)
}

// URLRule válida que el campo sea una URL válida
type URLRule struct{}

func (r *URLRule) Name() string {
	return "url"
}

func (r *URLRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("the field must be a string")
	}

	if !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "https://") {
		return fmt.Errorf("the field must be a valid URL")
	}

	return nil
}

// DateRule válida que el campo sea una fecha válida
type DateRule struct {
	Format string // Formato de fecha esperado (ej: "2006-01-02")
}

func (r *DateRule) Name() string {
	return "date"
}

func (r *DateRule) Validate(value interface{}, data map[string]interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("the field must be a string")
	}

	format := r.Format
	if format == "" {
		format = "2006-01-02" // Formato por defecto
	}

	// Validar la fecha usando 'time.Parse'
	if _, err := time.Parse(format, str); err != nil {
		return fmt.Errorf("the field must be a valid date")
	}

	return nil
}
