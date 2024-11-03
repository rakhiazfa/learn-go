package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rakhiazfa/vust-identity-service/api/dto/requests"
	"github.com/rakhiazfa/vust-identity-service/api/dto/responses"
	"github.com/rakhiazfa/vust-identity-service/pkg/utils"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type FileService struct {
	userContext *utils.UserContext
	serviceUrl  string
}

func NewFileService(userContext *utils.UserContext) *FileService {
	serviceUrl := viper.GetString("services.file_service")

	return &FileService{userContext, serviceUrl}
}

func (s *FileService) UploadFile(payload requests.UploadFileReq) (*responses.FileRes, error) {
	var buffer bytes.Buffer

	writer, err := utils.CreateFormFromStruct(&buffer, payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.serviceUrl+"/files", &buffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+s.userContext.GetAccessToken())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if closeErr := Body.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var fileAPIRes responses.FileAPIRes
	if err := json.Unmarshal(body, &fileAPIRes); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("upload failed with status: %s", res.Status)
	}

	return &fileAPIRes.File, err
}