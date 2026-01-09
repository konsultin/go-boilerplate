package service

import (
	"path"

	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/errk"
	logkOption "github.com/konsultin/logk/option"
)

func (s *Service) mustComposeFileResult(assetType dto.AssetType_Enum, fileName string) *dto.File {
	result, err := s.composeFileResult(assetType, fileName)
	if err != nil {
		panic(errk.Trace(err))
	}

	return result
}

func (s *Service) composeFileResult(assetType dto.AssetType_Enum, fileName string) (*dto.File, error) {
	if fileName == "" {
		return nil, nil
	}

	// Get asset path
	assetPath := getAssetPath(assetType, fileName)

	u, err := s.repo.GetDownloadFileUrl(assetPath)
	if err != nil {
		s.log.Error("Failed to resolve file url", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	// Set download url
	return &dto.File{
		FileName:  fileName,
		Url:       u,
		Signature: "",
	}, nil
}

func getAssetPath(assetType dto.AssetType_Enum, fileName string) string {
	assetDir := ""
	switch assetType {
	case dto.AssetType_USER_AVATAR:
		assetDir = "users/avatars"
	}
	// Join path
	return path.Join(assetDir, fileName)
}
