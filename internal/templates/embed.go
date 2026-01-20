// Package templates provides embedded template files for initializing .context/ directories.
package templates

import "embed"

//go:embed *.md
var FS embed.FS

// GetTemplate reads a template file by name from the embedded filesystem.
func GetTemplate(name string) ([]byte, error) {
	return FS.ReadFile(name)
}

// ListTemplates returns all available template file names.
func ListTemplates() ([]string, error) {
	entries, err := FS.ReadDir(".")
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
