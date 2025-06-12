package git

import (
	"fmt"
	"os"
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
			if len(line) < 4 {
				continue
			}

			// status puede estar en la primera o segunda columna
			status := strings.TrimSpace(line[:2])
			filename := strings.TrimSpace(line[3:])

			// Algunos nombres de archivo pueden contener espacios. Usa Fields para mayor precisión.
			// Nota: git status --porcelain v1 separa status y nombre con exactamente dos caracteres.
			if fields := strings.Fields(line); len(fields) >= 2 {
				filename = strings.Join(fields[1:], " ")
			}

			files = append(files, FileItem{
				Name:   filename,
				Status: status,
				Staged: stagedFiles[filename],
			})
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

// Fetch obtiene todos los cambios y ramas del repositorio remoto
func Fetch() error {
	output, err := exec.Command("git", "fetch", "-a").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error al hacer fetch: %v - %s", err, string(output))
	}
	return nil
}

// Pull obtiene e integra los cambios del repositorio remoto
func Pull() error {
	output, err := exec.Command("git", "pull").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error al hacer pull: %v - %s", err, string(output))
	}
	return nil
}

// CreateBranch crea una nueva rama
func CreateBranch(name string) error {
	cmd := exec.Command("git", "checkout", "-b", name)
	return cmd.Run()
}

// DeleteBranch elimina una rama
func DeleteBranch(name string) error {
	cmd := exec.Command("git", "branch", "-d", name)
	return cmd.Run()
}

// DiscardChanges descarta los cambios de un archivo
func DiscardChanges(filename string) error {
	// Verificar si el archivo está siendo rastreado por Git
	checkTracked := exec.Command("git", "ls-files", "--error-unmatch", filename)
	if err := checkTracked.Run(); err != nil {
		// El archivo NO está siendo rastreado, lo eliminamos
		fmt.Println("Archivo NO rastreado, intentando eliminar:", filename)
		if err := os.Remove(filename); err != nil {
			return fmt.Errorf("no se pudo eliminar archivo no rastreado: %w", err)
		}
		fmt.Println("Archivo eliminado correctamente")
		return nil
	}

	// El archivo está rastreado, se descartan los cambios
	fmt.Println("Archivo rastreado, descartando cambios con git checkout")
	discard := exec.Command("git", "checkout", "--", filename)
	if err := discard.Run(); err != nil {
		return fmt.Errorf("error al descartar cambios con git: %w", err)
	}

	fmt.Println("Cambios descartados correctamente")
	return nil
}

// HasRemoteChanges verifica si hay commits pendientes de pull desde el remoto
func HasRemoteChanges(branch string) (bool, error) {
	// Ejecuta git fetch para actualizar refs
	err := Fetch()
	if err != nil {
		return false, err
	}

	cmd := exec.Command("git", "rev-list", "--count", fmt.Sprintf("HEAD..origin/%s", branch))
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	count := strings.TrimSpace(string(output))
	return count != "0", nil
}
