package mModel

import "encoding/json"

type Course struct {
	ID              int             `json:"id"`
	Fullname        string          `json:"fullname"`
	ShortName       string          `json:"shortname"`
	Summary         string          `json:"summary"`
	Visible         bool            `json:"visible"`
	StartDate       int64           `json:"startdate"`
	EndDate         int64           `json:"enddate"`
	CourseImage     string          `json:"courseimage"`
	CourseImageType string          `json:"courseimagetype"`
	Category        string          `json:"coursecategory"`
	Sections        []CourseSection `json:"sections"`
}

type CourseSection struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	SectionNumber int            `json:"section"`
	Modules       []CourseModule `json:"modules"`
}

type CourseModule struct {
	ID                      int                     `json:"id"`
	Instance                int                     `json:"instance"`
	Description             string                  `json:"description"`
	URL                     string                  `json:"url"`
	Name                    string                  `json:"name"`
	ModIcon                 string                  `json:"modicon"`
	ModName                 string                  `json:"modname"`
	Contents                []CourseContent         `json:"contents"`
	CourseAssignmentContent CourseAssignmentContent `json:"assignment_content"`
}

type CourseContent struct {
	Type        string `json:"type"`
	FileName    string `json:"filename"`
	FileSize    int64  `json:"filesize"`
	FileURL     string `json:"fileurl"`
	TimeCreated int64  `json:"timecreated"`
}

type CourseAssignments struct {
	ID          int                   `json:"id"`
	Assignments []CourseModAssignment `json:"assignments"`
}
type CourseModAssignment struct {
	ComponentID         int          `json:"cmid"`
	ID                  int          `json:"id"`
	Intro               string       `json:"intro"`
	SubmissionStatement string       `json:"submissionstatement"`
	DueDate             int64        `json:"duedate"`
	IntroAttachment     []MoodleFile `json:"introattachments"`
}
type CourseModSubmission struct {
	GradeForDisplay       string       `json:"gradefordisplay"`
	GradedDate            int          `json:"gradeddate"`
	FeedbackComment       string       `json:"feedbackcomment"`
	SubmissionAttachments []MoodleFile `json:"submissionattachments"`
}
type CourseAssignmentContent struct {
	ID                    int          `json:"id"`
	Intro                 string       `json:"intro"`
	SubmissionStatement   string       `json:"submissionstatement"`
	DueDate               int64        `json:"duedate"`
	GradeForDisplay       string       `json:"gradefordisplay"`
	GradedDate            int          `json:"gradeddate"`
	FeedbackComment       string       `json:"feedbackcomment"`
	IntroAttachment       []MoodleFile `json:"introattachments"`
	SubmissionAttachments []MoodleFile `json:"submissionattachments"`
}
type MoodleFile struct {
	FileName     string `json:"filename"`
	FileSize     int64  `json:"filesize"`
	FileURL      string `json:"fileurl"`
	TimeModified int64  `json:"timemodified"`
}

func (c Course) MarshalJSON() ([]byte, error) {
	type Alias Course
	return json.Marshal(&struct {
		Sections interface{} `json:"sections,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(&c),
	})
}
