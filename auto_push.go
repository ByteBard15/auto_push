package main

import (
    "fmt"
    "os"
    "os/exec"
    "time"
    "path/filepath"
)

const lastPushFile = "last_push"
func main() {
	repoDir, _ := os.Getwd()
	pushFilePath := filepath.Join(repoDir, lastPushFile)

    needPush := shouldPush(pushFilePath)

	if needPush {
		if err := gitPush(repoDir); err != nil {
			fmt.Println("❌ Git push failed:", err)
			return
		}
		updateLastPush(pushFilePath)
		fmt.Println("✅ Repo pushed successfully at", time.Now().Format(time.RFC822))
	} else {
		fmt.Println("Last push was within 24 hours, skipping...")
	}
}

func shouldPush(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return true
	}
	lastTime, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return true
	}
	return time.Since(lastTime) >= 24*time.Hour
}

func updateLastPush(path string) {
	os.WriteFile(path, []byte(time.Now().Format(time.RFC3339)), 0644)
}

func gitPush(repoDir string) error {
	cmds := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", fmt.Sprintf("Auto backup: %s", time.Now().Format("2006-01-02 15:04:05"))},
		{"git", "push"},
	}

	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = repoDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
