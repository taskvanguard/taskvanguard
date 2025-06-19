package config

// "os"
// "path/filepath"

// func TemplatesDir() (string, error) {
// 	configDir, err := os.UserConfigDir()
// 	if err != nil {
// 		return "", err
// 	}
// 	return filepath.Join(configDir, "taskvanguard", "templates"), nil
// }

// func UserTemplatePath(filename string) (string, error) {
// 	dir, err := TemplatesDir()
// 	if err != nil {
// 		return "", err
// 	}
// 	return filepath.Join(dir, filename), nil
// }

// func TemplateExists(filename string) bool {
// 	path, err := UserTemplatePath(filename)
// 	if err != nil {
// 		return false
// 	}
// 	_, err = os.Stat(path)
// 	return err == nil
// }

// func EnsureTemplate(filename string, content string) error {
// 	path, err := UserTemplatePath(filename)
// 	if err != nil {
// 		return err
// 	}

// 	if _, err := os.Stat(path); err == nil {
// 		return nil // already exists
// 	}

// 	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
// 		return err
// 	}

// 	return os.WriteFile(path, []byte(content), 0644)
// }

// func ReadDefaultTemplate(filename string) (string, error) {
// 	data, err := Templates.ReadFile("templates/" + filename)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(data), nil
// }