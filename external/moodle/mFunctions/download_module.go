package mFunctions

import (
	"ScArium/external/moodle"
	"ScArium/external/moodle/mModel"
	"fmt"
	"os"
)

type downloadResourceModel struct {
	ID               int                 `json:"id"`
	Instance         int                 `json:"instance"`
	Name             string              `json:"name"`
	ModName          string              `json:"modName"`
	ContentFileNames []mModel.MoodleFile `json:"contentFileNames"`
}

func DownloadModule(module mModel.CourseModule, basePath string, mc *moodle.MoodleClient) error {
	switch module.ModName {
	case "label":
		return nil // nothing to download in labels
	case "resource":
		return downloadResourceModule(module, basePath, mc)
	case "url":
		return downloadUrlModule(module, basePath, mc)
	case "assign":
		return downloadAssignModule(module, basePath, mc)
	}
	return fmt.Errorf("Module not supported: %s", module.ModName)
}

func downloadResourceModule(m mModel.CourseModule, basePath string, mc *moodle.MoodleClient) error {
	path := fmt.Sprintf("%s/%d", basePath, m.ID)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	downloadModel := &downloadResourceModel{
		ID:               m.ID,
		Instance:         m.Instance,
		Name:             m.Name,
		ModName:          m.ModName,
		ContentFileNames: make([]mModel.MoodleFile, 0),
	}
	for _, file := range m.Contents {
		err := mc.DownloadFile(file.FileURL, file.FileSize, fmt.Sprintf("%s/%s", path, file.FileName))
		if err != nil {
			return err
		}
		newMoodleFile := mModel.MoodleFile{
			FileName:     file.FileName,
			FileSize:     file.FileSize,
			FileURL:      "/api/v1/moodle/",
			TimeModified: file.TimeCreated,
		}
		downloadModel.ContentFileNames = append(downloadModel.ContentFileNames, newMoodleFile)
	}
	//TODO: write data
	return nil

}
func downloadUrlModule(m mModel.CourseModule, path string, mc *moodle.MoodleClient) error {
	return nil
}
func downloadAssignModule(m mModel.CourseModule, path string, mc *moodle.MoodleClient) error {
	return nil
}
