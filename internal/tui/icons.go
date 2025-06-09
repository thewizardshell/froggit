package tui

import (
	"path/filepath"
	"strings"
)

func GetIconForFile(name string) string {
	ext := strings.ToLower(filepath.Ext(name))

	switch ext {
	case ".go":
		return "" // Go
	case ".js":
		return "" // JavaScript
	case ".ts":
		return "" // TypeScript
	case ".jsx":
		return "" // React JSX
	case ".tsx":
		return "" // React TSX
	case ".py":
		return "" // Python
	case ".java":
		return "" // Java
	case ".rb":
		return "" // Ruby
	case ".php":
		return "" // PHP
	case ".html", ".htm":
		return "" // HTML
	case ".css":
		return "" // CSS
	case ".json":
		return "" // JSON
	case ".md":
		return "" // Markdown
	case ".sh":
		return "" // Shell
	case ".yml", ".yaml":
		return "" // YAML
	case ".rs":
		return "" // Rust
	case ".cpp", ".cc", ".cxx", ".c++", ".h", ".hpp":
		return "" // C/C++
	case ".txt":
		return "" // Texto plano
	case ".lock":
		return "" // Archivo de lock
	case ".env":
		return "" // Archivo de entorno
	case ".svg", ".png", ".jpg", ".jpeg", ".gif", ".webp":
		return "" // Imágenes
	case ".exe":
		return "" // Ejecutable Windows
	case ".zip", ".tar", ".gz", ".rar":
		return "" // Archivo comprimido
	case ".log":
		return "" // Log
	default:
		return "" // Archivo genérico
	}
}
