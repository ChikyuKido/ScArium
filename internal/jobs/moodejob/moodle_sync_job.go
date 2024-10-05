package moodejob

import (
	"ScArium/common/config"
	"ScArium/external/moodle"
	"ScArium/external/moodle/mFunctions"
	"ScArium/external/moodle/mModel"
	"ScArium/internal/backend/database/repo"
	"ScArium/internal/jobs"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func NewMoodleSyncJob(accountId uint) error {
	err := jobs.CreateJob("MoodleSync", "For accountId "+strconv.Itoa(int(accountId)),
		map[string]string{
			"accountId": strconv.Itoa(int(accountId)),
		},
		syncMoodle,
	)
	if err != nil {
		return fmt.Errorf("failed to create moodle sync job: %w", err)
	}
	return nil
}

func syncMoodle(job jobs.Job) error {
	accountId, err := strconv.ParseUint(job.Params["accountId"], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to convert account id to uint: %w", err)
	}
	moodleAcc, err := repo.GetMoodleAccountById(uint(accountId))
	if err != nil {
		return fmt.Errorf("failed to get moodle account: %w", err)
	}
	mc, err := moodle.NewMoodleClient(moodleAcc.InstanceUrl, moodleAcc.Username, moodleAcc.Password)
	if err != nil {
		return fmt.Errorf("failed to create moodle client: %v", err)
	}
	courses, err := mFunctions.GetCourses(mc)
	if err != nil {
		return fmt.Errorf("failed to get courses: %v", err)
	}
	job.AppendLog(logrus.InfoLevel, fmt.Sprintf("Found %d courses", len(courses)))
	replaceImagesForCourses(job, courses)
	retrieveSections(job, courses, mc)
	retrieveAssignments(job, courses, mc)
	writeMetadata(job, courses, accountId)
	syncFiles(job, courses, accountId)
	return nil
}

func writeMetadata(job jobs.Job, courses []*mModel.Course, accountId uint64) {
	basePath := config.MOODLE_PATH + "/" + strconv.FormatUint(accountId, 10)
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to create metadata directory: %v", err))
		return
	}
	jsonData, err := json.Marshal(courses)
	if err != nil {
		job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to marshal courses: %v", err))
		return
	}
	err = os.WriteFile(basePath+"/courses.json", jsonData, 0666)
	if err != nil {
		job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to create metadata file: %v", err))
		return
	}
	for _, course := range courses {
		coursePath := basePath + "/" + strconv.Itoa(course.ID)
		err := os.MkdirAll(coursePath, os.ModePerm)
		if err != nil {
			job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to create metadata directory: %v", err))
			return
		}
		for _, section := range course.Sections {
			jsonData, err := json.Marshal(section)
			if err != nil {
				job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to marshal courses: %v", err))
				return
			}
			err = os.WriteFile(coursePath+"/section_"+strconv.Itoa(section.ID)+".json", jsonData, 0666)
			if err != nil {
				job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to create metadata file: %v", err))
				return
			}
		}
	}
}

func retrieveAssignments(job jobs.Job, courses []*mModel.Course, mc *moodle.MoodleClient) {
	courseAssignments, err := mFunctions.GetAssignments(mc, courses)
	if err != nil {
		job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to get assignments: %v", err))
	}
	for _, course := range courses {
		assignments, err := getAssignmentsForCourse(course.ID, courseAssignments)
		if err != nil {
			job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to get assignments: %v", err))
			continue
		}
		for i, sections := range course.Sections {
			for j, module := range sections.Modules {
				if module.ModName == "assign" {
					assignData, err := getAssignmentForId(module.Instance, assignments)
					if err != nil {
						job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to get assignment data: %v", err))
						continue
					}
					submissionData, err := mFunctions.GetSubmissionData(mc, module.Instance)
					if err != nil {
						job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to get assignment submission data: %v", err))
						continue
					}
					assignContent := mModel.CourseAssignmentContent{
						ComponentID:           assignData.ComponentID,
						ID:                    assignData.ID,
						Intro:                 assignData.Intro,
						SubmissionStatement:   assignData.SubmissionStatement,
						DueDate:               assignData.DueDate,
						GradeForDisplay:       submissionData.GradeForDisplay,
						GradedDate:            submissionData.GradedDate,
						FeedbackComment:       submissionData.FeedbackComment,
						IntroAttachment:       assignData.IntroAttachment,
						SubmissionAttachments: submissionData.SubmissionAttachments,
					}
					course.Sections[i].Modules[j].CourseAssignmentContent = assignContent
				}
			}
		}
	}
}

func getAssignmentForId(id int, assignments *mModel.CourseAssignments) (mModel.CourseModAssignment, error) {
	for _, assign := range assignments.Assignments {
		if assign.ID == id {
			return assign, nil
		}
	}
	return mModel.CourseModAssignment{}, fmt.Errorf("failed to retrieve assignment for id: %d", id)
}

func getAssignmentsForCourse(id int, courseAssignments []*mModel.CourseAssignments) (*mModel.CourseAssignments, error) {
	for _, assignemnts := range courseAssignments {
		if assignemnts.ID == id {
			return assignemnts, nil
		}
	}
	return nil, fmt.Errorf("failed to retrieve assignments for course with id: %d", id)
}

func retrieveSections(job jobs.Job, courses []*mModel.Course, mc *moodle.MoodleClient) {
	for _, course := range courses {
		err := mFunctions.GetSections(mc, course)
		if err != nil {
			job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("failed to get sections for course: %s: %v", course.ShortName, err))
		}
		job.AppendLog(logrus.InfoLevel, fmt.Sprintf("Found sections %d for course: %s", len(course.Sections), course.ShortName))
	}
}

func replaceImagesForCourses(job jobs.Job, courses []*mModel.Course) {
	for i := range courses {
		if strings.Contains(courses[i].CourseImage, "http") {
			file, err := os.ReadFile("static/imgs/defaultMoodleCourse.svg")
			if err != nil {
				job.AppendLog(logrus.ErrorLevel, "Could not load default svg")
			}
			job.AppendLog(logrus.WarnLevel, fmt.Sprintf("%s does not have a image load default", courses[i].ShortName))
			courses[i].CourseImage = string(file)
		}
	}
}
