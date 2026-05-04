package utils

import (
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

const (
	DocumenFileDir = "document"
)

func SaveFileToPath(file *multipart.FileHeader, folderName string, ctx fiber.Ctx) (*string, error) {
	randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	path := fmt.Sprintf("/uploads/%s/%s", folderName, strconv.FormatInt(int64(randomizer.Int()), 32)+file.Filename)

	return saveFile(file, path, ctx)
}

func SaveFileToCustomPath(file *multipart.FileHeader, folderPath string, ctx fiber.Ctx) (*string, error) {
	randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	cleanFolderPath := strings.Trim(folderPath, "/")
	path := fmt.Sprintf("/uploads/%s/%s", cleanFolderPath, strconv.FormatInt(int64(randomizer.Int()), 32)+file.Filename)

	return saveFile(file, path, ctx)
}

func saveFile(file *multipart.FileHeader, path string, ctx fiber.Ctx) (*string, error) {
	fullPath := "." + path

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return nil, err
	}

	err := ctx.SaveFile(file, fullPath)

	if err != nil {
		return nil, err
	}

	return &path, nil
}

func RemoveFileFromPath(path string) error {
	// Check if the file exists
	if _, err := os.Stat("." + path); os.IsNotExist(err) {
		return fiber.NewError(404, "File cannot be found")
	}

	// Remove the file
	if err := os.Remove("." + path); err != nil {
		return fiber.NewError(500, "Failed to remove file")
	}

	return nil
}
