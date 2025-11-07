package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "time"
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
    addCmd := exec.Command("git", "add", ".")
    addCmd.Dir = repoDir
    addCmd.Stdout = os.Stdout
    addCmd.Stderr = os.Stderr
    if err := addCmd.Run(); err != nil {
        return err
    }

    commitCmd := exec.Command("git", "commit", "-m", fmt.Sprintf("Auto backup: %s", time.Now().Format("2006-01-02 15:04:05")))
    commitCmd.Dir = repoDir
    commitCmd.Stdout = os.Stdout
    commitCmd.Stderr = os.Stderr
    err := commitCmd.Run()
    if err != nil {
        // Check if it's just "nothing to commit" — in which case we still continue
        exitErr, ok := err.(*exec.ExitError)
        if ok && exitErr.ExitCode() == 1 {
            fmt.Println("ℹ️ Nothing new to commit, proceeding with push anyway...")
        } else {
            return err
        }
    }

    pushCmd := exec.Command("git", "push")
    pushCmd.Dir = repoDir
    pushCmd.Stdout = os.Stdout
    pushCmd.Stderr = os.Stderr
    if err := pushCmd.Run(); err != nil {
        return err
    }

    return nil
}

