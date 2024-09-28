package mFunctions

import (
	"ScArium/external/moodle"
	"ScArium/external/moodle/mModel"
	"encoding/json"
)

func GetCourses(client *moodle.MoodleClient) ([]*mModel.Course, error) {
	body, err := client.MakeWebserviceRequest("core_course_get_enrolled_courses_by_timeline_classification", map[string]string{"classification": "all"})

	if err != nil {
		return nil, err
	}
	var coursesResp struct {
		Courses []mModel.Course `json:"courses"`
		Error   string          `json:"error"`
	}
	if err := json.Unmarshal(body, &coursesResp); err != nil {
		return nil, err
	}
	courses := make([]*mModel.Course, len(coursesResp.Courses))
	for i := range coursesResp.Courses {
		courses[i] = &coursesResp.Courses[i]
	}

	return courses, nil
}
