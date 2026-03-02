package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// IsGitRepo checks if the given directory is a git repository
func IsGitRepo(dir string) bool {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// GetDiff gets all uncommitted changes (staged + unstaged)
func GetDiff(dir string) (string, error) {
	// First check if it's a git repo
	if !IsGitRepo(dir) {
		return "", fmt.Errorf("not a git repository")
	}

	// Get staged changes
	stagedCmd := exec.Command("git", "-C", dir, "diff", "--cached")
	var stagedOut bytes.Buffer
	stagedCmd.Stdout = &stagedOut
	if err := stagedCmd.Run(); err != nil {
		// It's okay if there are no staged changes
	}

	// Get unstaged changes
	unstagedCmd := exec.Command("git", "-C", dir, "diff")
	var unstagedOut bytes.Buffer
	unstagedCmd.Stdout = &unstagedOut
	if err := unstagedCmd.Run(); err != nil {
		// It's okay if there are no unstaged changes
	}

	diff := stagedOut.String() + unstagedOut.String()

	if strings.TrimSpace(diff) == "" {
		// Try to get untracked files
		untrackedCmd := exec.Command("git", "-C", dir, "ls-files", "--others", "--exclude-standard")
		var untrackedOut bytes.Buffer
		untrackedCmd.Stdout = &untrackedOut
		if err := untrackedCmd.Run(); err == nil {
			untracked := untrackedOut.String()
			if strings.TrimSpace(untracked) != "" {
				diff = fmt.Sprintf("[Untracked files:\n%s\n]", untracked)
			}
		}
	}

	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("no changes detected")
	}

	return diff, nil
}

// GetStatus returns the git status output
func GetStatus(dir string) (string, error) {
	if !IsGitRepo(dir) {
		return "", fmt.Errorf("not a git repository")
	}

	cmd := exec.Command("git", "-C", dir, "status", "--short")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}

// Commit executes git commit with the given message
func Commit(dir, message string) error {
	// First, stage all changes to ensure everything gets committed
	addCmd := exec.Command("git", "-C", dir, "add", "-A")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to stage changes: %w", err)
	}

	// Then commit
	cmd := exec.Command("git", "-C", dir, "commit", "-m", message)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if errMsg := stderr.String(); errMsg != "" {
			return fmt.Errorf("git commit failed: %s", strings.TrimSpace(errMsg))
		}
		return fmt.Errorf("git commit failed: %w", err)
	}
	return nil
}
