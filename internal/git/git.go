package git

import (
	"os/exec"
	"strings"
)

// FileItem representa un archivo en el repositorio
type FileItem struct {
	Name     string
	Status   string
	Staged   bool
	Selected bool
}

// IsGitRepository verifica si el directorio actual es un repositorio Git
func IsGitRepository() bool {
	_, err := exec.Command("git", "rev-parse", "--git-dir").Output()
	return err == nil
}

// InitRepository inicializa un nuevo repositorio Git
func InitRepository() error {
	cmd := exec.Command("git", "init")
	return cmd.Run()
}

// GetModifiedFiles obtiene la lista de archivos modificados
func GetModifiedFiles() ([]FileItem, error) {
	// Obtener archivos staged
	stagedCmd := exec.Command("git", "diff", "--cached", "--name-status")
	stagedOutput, _ := stagedCmd.Output()
	stagedFiles := make(map[string]bool)

	if len(stagedOutput) > 0 {
		lines := strings.Split(strings.TrimSpace(string(stagedOutput)), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					stagedFiles[parts[1]] = true
				}
			}
		}
	}

	// Obtener todos los archivos modificados
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var files []FileItem
	if len(output) > 0 {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			if line != "" && len(line) >= 3 {
				status := strings.TrimSpace(line[:2])
				filename := line[3:]

				files = append(files, FileItem{
					Name:   filename,
					Status: status,
					Staged: stagedFiles[filename],
				})
			}
		}
	}

	return files, nil
}

// GetBranches obtiene la lista de ramas y la rama actual
func GetBranches() ([]string, string) {
	cmd := exec.Command("git", "branch")
	output, err := cmd.Output()
	if err != nil {
		return []string{}, ""
	}

	var branches []string
	var current string

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "* ") {
			current = line[2:]
			branches = append(branches, current)
		} else {
			branches = append(branches, line)
		}
	}

	return branches, current
}

// GetRemotes obtiene la lista de remotes configurados
func GetRemotes() ([]string, error) {
	cmd := exec.Command("git", "remote", "-v")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var remotes []string
	if len(output) > 0 {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		seen := make(map[string]bool)
		for _, line := range lines {
			if line != "" {
				parts := strings.Fields(line)
				if len(parts) >= 2 && !seen[parts[0]] {
					remotes = append(remotes, parts[0]+" -> "+parts[1])
					seen[parts[0]] = true
				}
			}
		}
	}

	return remotes, nil
}

// Add añade un archivo al staging area
func Add(filename string) error {
	cmd := exec.Command("git", "add", filename)
	return cmd.Run()
}

// Reset quita un archivo del staging area
func Reset(filename string) error {
	cmd := exec.Command("git", "reset", "HEAD", filename)
	return cmd.Run()
}

// Commit realiza un commit con el mensaje especificado
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}

// Push sube los cambios al repositorio remoto
func Push() error {
	cmd := exec.Command("git", "push")
	return cmd.Run()
}

// Checkout cambia a la rama especificada
func Checkout(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	return cmd.Run()
}

// AddRemote añade un remote al repositorio
func AddRemote(name, url string) error {
	cmd := exec.Command("git", "remote", "add", name, url)
	return cmd.Run()
}

// RemoveRemote elimina un remote del repositorio
func RemoveRemote(name string) error {
	cmd := exec.Command("git", "remote", "remove", name)
	return cmd.Run()
}
