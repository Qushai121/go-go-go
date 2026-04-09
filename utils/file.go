package utils

import (
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

const (
	DocumenFileDir = "document"
)

func SaveFileToPath(file *multipart.FileHeader, folderName string, ctx fiber.Ctx) (*string, error) {
	randomizer := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	path := fmt.Sprintf("/uploads/%s/%s", folderName, strconv.FormatInt(int64(randomizer.Int()), 32)+file.Filename)

	err := ctx.SaveFile(file, "."+path)

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
