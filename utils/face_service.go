package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type FaceServiceResult struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

func VerifyFaceImage(file *multipart.FileHeader) (*FaceServiceResult, error) {
	return postFaceServiceMultipart("/api/only-verify-face", file, nil)
}

func VerifyFaceByNIK(file *multipart.FileHeader, nik string) (*FaceServiceResult, error) {
	return postFaceServiceMultipart("/api/verify-face", file, map[string]string{
		"nik": nik,
	})
}

func FaceSimilarityPercentage(result *FaceServiceResult) string {
	if result == nil || result.Data == nil {
		return ""
	}

	value, ok := result.Data["presentase_kemiripan"]
	if !ok || value == nil {
		return ""
	}

	return fmt.Sprint(value)
}

func postFaceServiceMultipart(path string, file *multipart.FileHeader, fields map[string]string) (*FaceServiceResult, error) {
	if file == nil {
		return nil, fmt.Errorf("image file is required")
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	part, err := writer.CreateFormFile("image", file.Filename)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, src); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	baseURL := strings.TrimRight(os.Getenv("FACE_SERVICE_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:7708"
	}

	client := &http.Client{Timeout: faceServiceTimeout()}
	req, err := http.NewRequest(http.MethodPost, baseURL+path, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("face service is unavailable: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result FaceServiceResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("invalid face service response: %w", err)
	}

	if result.Status == 0 {
		result.Status = resp.StatusCode
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 || !result.Success {
		if strings.TrimSpace(result.Message) == "" {
			result.Message = "face verification failed"
		}
		return &result, fmt.Errorf(result.Message)
	}

	return &result, nil
}

func faceServiceTimeout() time.Duration {
	raw := strings.TrimSpace(os.Getenv("FACE_SERVICE_TIMEOUT_SECONDS"))
	if raw == "" {
		return 30 * time.Second
	}

	duration, err := time.ParseDuration(raw + "s")
	if err != nil || duration <= 0 {
		return 30 * time.Second
	}

	return duration
}
