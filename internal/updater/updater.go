package updater

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/inconshreveable/go-update"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const repo = "thewizardshell/froggit"

func CheckAndUpdate(version string) {
	v, err := semver.ParseTolerant(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid version format %q: %v\n", version, err)
		return
	}

	fmt.Println("Checking for updates...")

	// DetectLatest already picks the asset matching current OS/arch
	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		return
	}

	if !found {
		fmt.Println("No releases found.")
		return
	}

	if !latest.Version.GT(v) {
		fmt.Printf("You are already on the latest version (%s).\n", v)
		return
	}

	fmt.Printf("New version available: %s (current: %s)\n", latest.Version, v)
	if latest.ReleaseNotes != "" {
		fmt.Printf("Release notes:\n%s\n\n", latest.ReleaseNotes)
	}

	fmt.Print("Do you want to update now? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		fmt.Println("Update cancelled.")
		return
	}

	if err := applyUpdate(latest.AssetURL); err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		return
	}

	fmt.Printf("Successfully updated to version %s. Please restart froggit.\n", latest.Version)
}

// applyUpdate downloads the zip from assetURL, extracts the first binary,
// and replaces the current executable — regardless of the binary name inside the zip.
func applyUpdate(assetURL string) error {
	fmt.Println("Downloading...")

	resp, err := http.Get(assetURL)
	if err != nil {
		return fmt.Errorf("download error: %w", err)
	}
	defer resp.Body.Close()

	tmp, err := os.CreateTemp("", "froggit-update-*.zip")
	if err != nil {
		return fmt.Errorf("temp file error: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	tmp.Close()

	binary, err := extractBinaryFromZip(tmp.Name())
	if err != nil {
		return err
	}
	defer binary.Close()

	return update.Apply(binary, update.Options{})
}

// extractBinaryFromZip opens a zip and returns a reader for the first non-directory file.
func extractBinaryFromZip(zipPath string) (io.ReadCloser, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("zip open error: %w", err)
	}

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			r.Close()
			return nil, fmt.Errorf("zip entry error: %w", err)
		}
		// Wrap so closing also closes the zip reader
		return &zipEntry{rc: rc, zr: r}, nil
	}

	r.Close()
	return nil, fmt.Errorf("no binary found in zip")
}

type zipEntry struct {
	rc io.ReadCloser
	zr *zip.ReadCloser
}

func (z *zipEntry) Read(p []byte) (int, error) { return z.rc.Read(p) }
func (z *zipEntry) Close() error {
	z.rc.Close()
	return z.zr.Close()
}
