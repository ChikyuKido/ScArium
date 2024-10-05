package mFunctions

import (
	"ScArium/external/moodle"
	"ScArium/external/moodle/mModel"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func GetAssignments(client *moodle.MoodleClient, courses []*mModel.Course) ([]*mModel.CourseAssignments, error) {
	var ids []string
	for _, course := range courses {
		ids = append(ids, strconv.Itoa(course.ID))
	}
	var params = map[string]string{}
	for i := 0; i < len(ids); i++ {
		params[fmt.Sprintf("courseids[%d]", i)] = ids[i]
	}
	body, err := client.MakeWebserviceRequest("mod_assign_get_assignments", params)
	if err != nil {
		return nil, err
	}

	var coursesResp struct {
		Courses []*mModel.CourseAssignments `json:"courses"`
		Error   string                      `json:"error"`
	}
	if err := json.Unmarshal(body, &coursesResp); err != nil {
		return nil, err
	}
	if coursesResp.Error != "" {
		return nil, errors.New(coursesResp.Error)
	}

	return coursesResp.Courses, nil
}
func GetSubmissionData(client *moodle.MoodleClient, id int) (mModel.CourseModSubmission, error) {
	var body, err = client.MakeWebserviceRequest("mod_assign_get_submission_status", map[string]string{"assignid": strconv.Itoa(id)})
	if err != nil {
		return mModel.CourseModSubmission{}, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return mModel.CourseModSubmission{}, err
	}

	submissionAttachments, err := getSubmissionFiles(result)
	if err != nil {
		return mModel.CourseModSubmission{}, err
	}

	submissionData := mModel.CourseModSubmission{}
	submissionData.SubmissionAttachments = submissionAttachments
	err = fillFeedbackData(&submissionData, result)
	if err != nil {
		return mModel.CourseModSubmission{}, err
	}

	return submissionData, nil
}
func getSubmissionFiles(json map[string]interface{}) ([]mModel.MoodleFile, error) {
	var submissionMoodleFiles []mModel.MoodleFile
	if value, ok := json["lastattempt"]; ok {
		lastAttempt := value.(map[string]interface{})
		if value, ok := lastAttempt["submission"]; ok {
			submission := value.(map[string]interface{})
			if value, ok := submission["plugins"]; ok {
				plugins := value.([]interface{})
				for _, plugin := range plugins {
					pluginMap := plugin.(map[string]interface{})
					if pluginMap["type"] == "file" {
						fileAreas := pluginMap["fileareas"].([]interface{})
						for _, fileArea := range fileAreas {
							fileAreaMap := fileArea.(map[string]interface{})
							files := fileAreaMap["files"].([]interface{})
							for _, file := range files {
								fileMap := file.(map[string]interface{})
								moodleFile := mModel.MoodleFile{
									FileName:     fileMap["filename"].(string),
									FileSize:     int64(fileMap["filesize"].(float64)),
									FileURL:      fileMap["fileurl"].(string),
									TimeModified: int64(fileMap["timemodified"].(float64)),
								}
								submissionMoodleFiles = append(submissionMoodleFiles, moodleFile)
							}
						}
					}
				}
			}
		}
	}
	return submissionMoodleFiles, nil
}

func fillFeedbackData(submissionData *mModel.CourseModSubmission, json map[string]interface{}) error {
	if feedback, ok := json["feedback"].(map[string]interface{}); ok {

		if gradeForDisplay, ok := feedback["gradefordisplay"].(string); ok {
			submissionData.GradeForDisplay = gradeForDisplay
		} else {
			return errors.New("gradefordisplay not found")
		}

		if gradedDate, ok := feedback["gradeddate"].(float64); ok {
			submissionData.GradedDate = int(gradedDate)
		} else {
			return errors.New("gradeddate not found")
		}

		if plugins, ok := feedback["plugins"].([]interface{}); ok {
			for _, plugin := range plugins {
				pluginMap := plugin.(map[string]interface{})
				if pluginMap["type"] == "comments" {
					if editorFields, ok := pluginMap["editorfields"].([]interface{}); ok {
						for _, editorField := range editorFields {
							editorFieldMap := editorField.(map[string]interface{})
							if editorFieldMap["name"] == "comments" {
								if feedbackComment, ok := editorFieldMap["text"].(string); ok {
									submissionData.FeedbackComment = feedbackComment
									break
								}
							}
						}
					}
				}
			}
		}
	} else {
		return errors.New("feedback not found")
	}
	return nil
}
