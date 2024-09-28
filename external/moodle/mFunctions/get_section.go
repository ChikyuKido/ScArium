package mFunctions

import (
	"ScArium/external/moodle"
	"ScArium/external/moodle/mModel"
	"encoding/json"
	"fmt"
	"strconv"
)

func GetSections(client *moodle.MoodleClient, course *mModel.Course) error {
	var body, err = client.MakeWebserviceRequest("core_course_get_contents", map[string]string{"courseid": strconv.Itoa(course.ID)})

	if err != nil {
		return err
	}
	//When the json starts with a bracket then it's an error because a valid start with a [
	if body[0] == '{' {
		return fmt.Errorf("error in response json: %v", string(body))
	}
	var sections []mModel.CourseSection
	if err := json.Unmarshal(body, &sections); err != nil {
		return err
	}
	course.Sections = sections
	return nil
}
