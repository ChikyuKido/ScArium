package repo

import (
	"ScArium/internal/backend/database"
	"ScArium/internal/backend/database/entity"
	"log"
	"time"
)

func InsertCourse(course entity.MoodleCourse) error {
	startTime := time.Now() // Start the timer

	err := database.GetDB().Create(&course).Error /*.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "course_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"fullname", "shortname", "summary", "visible", "start_date", "end_date", "course_image", "course_image_type", "course_category"}),
	})*/

	elapsedTime := time.Since(startTime)            // Calculate elapsed time
	log.Printf("InsertCourse took %s", elapsedTime) // Log the time taken

	return err
}

func InsertSections(sections []entity.MoodleCourseSection) error {
	startTime := time.Now() // Start the timer

	err := database.GetDB().Create(&sections).Error /*.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "section_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "SectionNumber"}),
	})*/

	elapsedTime := time.Since(startTime)              // Calculate elapsed time
	log.Printf("InsertSections took %s", elapsedTime) // Log the time taken

	return err
}

func InsertResource(modules []entity.MoodleResource) error {
	startTime := time.Now() // Start the timer

	err := database.GetDB().Create(&modules).Error /*(.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "resource_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"instance", "description", "url", "name", "modicon", "modname", "resource_content", "assign_intro_resource"}),
	})*/

	elapsedTime := time.Since(startTime)              // Calculate elapsed time
	log.Printf("InsertResource took %s", elapsedTime) // Log the time taken

	return err
}

func InsertResourceAssign(modules []entity.MoodleAssignSubmissionResource) error {
	startTime := time.Now() // Start the timer

	err := database.GetDB().Create(&modules).Error /*.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "assign_submission_resource_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"grade_for_display", "graded_data", "feedback_comment", "submission_attachments", "account_id", "moodle_resource_id"}),
	})*/

	elapsedTime := time.Since(startTime)                    // Calculate elapsed time
	log.Printf("InsertResourceAssign took %s", elapsedTime) // Log the time taken

	return err
}
