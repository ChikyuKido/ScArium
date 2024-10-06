package moodejob

import (
	"ScArium/external/moodle"
	"ScArium/external/moodle/mFunctions"
	"ScArium/external/moodle/mModel"
	"ScArium/internal/backend/database/entity"
	"ScArium/internal/backend/database/repo"
	"ScArium/internal/backend/helper"
	"ScArium/internal/jobs"
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
	syncFiles(job, courses, accountId, mc)
	return nil
}

func writeMetadata(job jobs.Job, courses []*mModel.Course, accountId uint64) {
	for _, course := range courses {
		databaseCourse := entity.MoodleCourse{
			CourseID:        uint(course.ID),
			Fullname:        course.Fullname,
			ShortName:       course.ShortName,
			Summary:         course.Summary,
			Visible:         course.Visible,
			StartDate:       course.StartDate,
			EndDate:         course.EndDate,
			CourseImage:     course.CourseImage,
			CourseImageType: course.CourseImageType,
			Category:        course.Category,
		}
		err := repo.InsertCourse(databaseCourse)
		if err != nil {
			job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("failed to insert course with id %d to database", course.ID))
		}
		databaseSections := make([]entity.MoodleCourseSection, len(course.Sections))
		for i, section := range course.Sections {
			databaseSections[i].SectionID = uint(course.Sections[i].ID)
			databaseSections[i].SectionNumber = course.Sections[i].SectionNumber
			databaseSections[i].Name = course.Sections[i].Name
			databaseSections[i].MoodleCourseId = uint(course.ID)
			databaseResources := make([]entity.MoodleResource, len(section.Modules))
			databaseAssignSubmissionResources := make([]entity.MoodleAssignSubmissionResource, 0)
			for i, module := range section.Modules {
				intro := entity.MoodleAssignIntroResource{
					Intro:               module.CourseAssignmentContent.Intro,
					SubmissionStatement: module.CourseAssignmentContent.SubmissionStatement,
					IntroAttachment:     helper.ConvertStructToJsonOr(module.CourseAssignmentContent.IntroAttachment, "[]"),
				}

				databaseResources[i].ResourceID = uint(module.ID)
				databaseResources[i].Instance = uint(module.Instance)
				databaseResources[i].ModName = module.ModName
				databaseResources[i].Name = module.Name
				databaseResources[i].MoodleCourseSectionId = uint(section.ID)
				databaseResources[i].Description = module.Description
				databaseResources[i].ModIcon = module.ModIcon
				databaseResources[i].URL = module.URL
				databaseResources[i].AssignIntroResource = helper.ConvertStructToJsonOr(intro, "{}")
				databaseResources[i].ResourceContent = helper.ConvertStructToJsonOr(module.Contents, "[]")

				if module.ModName == "assign" {
					submissionResource := entity.MoodleAssignSubmissionResource{
						GradeForDisplay:       module.CourseAssignmentContent.GradeForDisplay,
						GradedDate:            module.CourseAssignmentContent.GradedDate,
						FeedbackComment:       module.CourseAssignmentContent.FeedbackComment,
						SubmissionAttachments: helper.ConvertStructToJsonOr(module.CourseAssignmentContent.SubmissionAttachments, "[]"),
						AccountId:             uint(accountId),
						MoodleResourceID:      uint(module.ID),
					}
					databaseAssignSubmissionResources = append(databaseAssignSubmissionResources, submissionResource)
				}
			}
			err = repo.InsertResource(databaseResources)
			if err != nil {
				job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("failed to insert resources for course with id %d to database", course.ID))
			}
			err = repo.InsertResourceAssign(databaseAssignSubmissionResources)
			if err != nil {
				job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("failed to insert submission resources for course with id %d to database", course.ID))
			}
		}
		err = repo.InsertSections(databaseSections)
		if err != nil {
			job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("failed to insert sections for course with id %d to database", course.ID))
		}

		fmt.Println(course.ShortName)
	}
}

/*
func writeMetadata(job jobs.Job, courses []*mModel.Course, accountId uint64) {
	job.AppendLog(logrus.InfoLevel, "Writing metadata for  courses")
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
	job.AppendLog(logrus.InfoLevel, "Metadata written successfully")
}*/

func retrieveAssignments(job jobs.Job, courses []*mModel.Course, mc *moodle.MoodleClient) {
	job.AppendLog(logrus.InfoLevel, "Retrieve assignments for courses")
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
	job.AppendLog(logrus.InfoLevel, "Finish retrieve assignments for courses")
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
	job.AppendLog(logrus.InfoLevel, "Retrieving sections for courses")
	for _, course := range courses {
		err := mFunctions.GetSections(mc, course)
		if err != nil {
			job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("failed to get sections for course: %s: %v", course.ShortName, err))
		}
	}
	job.AppendLog(logrus.InfoLevel, "Finish retrieving sections for courses")
}

func replaceImagesForCourses(job jobs.Job, courses []*mModel.Course) {
	job.AppendLog(logrus.InfoLevel, "Replace images for courses")
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
	job.AppendLog(logrus.InfoLevel, "Finish replacing images for courses")
}
