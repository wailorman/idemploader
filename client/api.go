package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/sendgrid/rest"
)

// IdemploaderAPI _
type IdemploaderAPI struct {
	baseURL     string
	accessToken string
}

// File _
type File struct {
	URL      string `json:"url"`
	Size     int    `json:"size"`
	Checksum string `json:"checksum"`
	MimeType string `json:"mime_type"`
}

// Config _
type Config struct {
	Host        string
	AccessToken string
}

// New _
func New(cfg Config) *IdemploaderAPI {
	return &IdemploaderAPI{
		baseURL:     cfg.Host,
		accessToken: cfg.AccessToken,
	}
}

// fileField _
const fileField = "file"

// AccessTokenHeaderName _
const AccessTokenHeaderName = "X-Access-Token"

// UploadFromFile _
func (ia *IdemploaderAPI) UploadFromFile(filepath string) (*File, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, errors.New("Can't open token's database file: " + err.Error())
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fileField, filepath)

	if err != nil {
		return nil, errors.New("Can't put token file to payload for API: " + err.Error())
	}

	io.Copy(part, file)

	writer.Close()

	req, err := http.NewRequest("POST", ia.baseURL+"/v1/files", body)

	if err != nil {
		return nil, errors.New("Failed to build request: " + err.Error())
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set(AccessTokenHeaderName, ia.accessToken)
	req.Header.Set("Accept", "application/json")

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, errors.New("Failed to make an API request: " + err.Error())
	}

	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, errors.New("Can't read response body: " + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("API error: " + string(responseBody))
	}

	responseUploader := &File{}

	err = json.Unmarshal(responseBody, &responseUploader)

	if err != nil {
		return nil, errors.New("Can't unmarshal response json: " + err.Error())
	}

	return responseUploader, nil
}

// GetFileInfo _
func (ia *IdemploaderAPI) GetFileInfo(checksum string) (*File, error) {
	request := rest.Request{
		Method:  rest.Get,
		BaseURL: ia.baseURL + "/v1/files/" + checksum,
		Headers: map[string]string{
			AccessTokenHeaderName: ia.accessToken,
		},
	}

	response, err := rest.Send(request)

	if err != nil {
		return nil, fmt.Errorf("Request error: %s", err.Error())
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP Error: %d: %s", response.StatusCode, response.Body)
	}

	fileInfo := &File{}

	err = json.Unmarshal([]byte(response.Body), fileInfo)

	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %s", err.Error())
	}

	return fileInfo, nil
}

// DownloadFileToPath _
func DownloadFileToPath(filepath, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(filepath)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

// DownloadFile _
func DownloadFile(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
