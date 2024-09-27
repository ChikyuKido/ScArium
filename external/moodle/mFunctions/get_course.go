package mFunctions

import (
	"ScArium/common/log"
	"ScArium/external/moodle"
	"ScArium/external/moodle/mModel"
	"encoding/json"
	"os"
	"strings"
)

func GetCourses(client *moodle.MoodleClient) ([]mModel.Course, error) {
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
	log.E.Info("Found ", len(coursesResp.Courses), " Course")
	for i := range coursesResp.Courses {
		if strings.Contains(coursesResp.Courses[i].CourseImage, "http") {
			file, err := os.ReadFile("static/imgs/defaultMoodleCourse.svg")
			if err != nil {
				log.E.Error("Could not load default svg")
			}
			log.E.Infof("%s does not have a image load default", coursesResp.Courses[i].ShortName)
			coursesResp.Courses[i].CourseImage = string(file)
		}
	}
	return coursesResp.Courses, nil
}
