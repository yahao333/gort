package backup

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

type BackupManager struct {
    backupDir string
    logger    *Logger
}

type BackupMetadata struct {
    Timestamp   time.Time `json:"timestamp"`
    Environment string    `json:"environment"`
    Version     string    `json:"version"`
    Type        string    `json:"type"`
}

func NewBackupManager(backupDir string, logger *Logger) *BackupManager {
    return &BackupManager{
        backupDir: backupDir,
        logger:    logger,
    }
}

func (bm *BackupManager) CreateBackup(env string, sourceDir string) (string, error) {
    timestamp := time.Now().Format("20060102-150405")
    backupFile := filepath.Join(bm.backupDir, fmt.Sprintf("%s-%s.zip", env, timestamp))

    if err := os.MkdirAll(bm.backupDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create backup directory: %w", err)
    }

    zipfile, err := os.Create(backupFile)
    if err != nil {
        return "", fmt.Errorf("failed to create backup file: %w", err)
    }
    defer zipfile.Close()

    archive := zip.NewWriter(zipfile)
    defer archive.Close()

    err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        relPath, err := filepath.Rel(sourceDir, path)
        if err != nil {
            return err
        }

        zipFile, err := archive.Create(relPath)
        if err != nil {
            return err
        }

        fsFile, err := os.Open(path)
        if err != nil {
            return err
        }
        defer fsFile.Close()

        _, err = io.Copy(zipFile, fsFile)
        return err
    })

    if err != nil {
        return "", fmt.Errorf("failed to create backup: %w", err)
    }

    bm.logger.Infof("Created backup: %s", backupFile)
    return backupFile, nil
}

func (bm *BackupManager) RestoreBackup(backupFile string, targetDir string) error {
    reader, err := zip.OpenReader(backupFile)
    if err != nil {
        return fmt.Errorf("failed to open backup file: %w", err)
    }
    defer reader.Close()

    for _, file := range reader.File {
        path := filepath.Join(targetDir, file.Name)

        if file.FileInfo().IsDir() {
            os.MkdirAll(path, file.Mode())
            continue
        }

        if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
            return fmt.Errorf("failed to create directory: %w