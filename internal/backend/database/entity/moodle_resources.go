package entity

type MoodleCourse struct {
	ID                  uint                  `json:"id" gorm:"primarykey"`
	CourseID            uint                  `json:"course_id"`
	Fullname            string                `json:"fullname"`
	ShortName           string                `json:"shortname"`
	Summary             string                `json:"summary"`
	Visible             bool                  `json:"visible"`
	StartDate           int64                 `json:"startdate"`
	EndDate             int64                 `json:"enddate"`
	CourseImage         string                `json:"courseimage"`
	CourseImageType     string                `json:"courseimagetype"`
	Category            string                `json:"coursecategory"`
	MoodleCourseSection []MoodleCourseSection `gorm:"foreignKey:MoodleCourseId"`
}

type MoodleCourseSection struct {
	ID             uint             `json:"id" gorm:"primarykey"`
	SectionID      uint             `json:"section_id"`
	Name           string           `json:"name"`
	SectionNumber  int              `json:"section"`
	MoodleCourseId uint             `json:"moodlecourseid" gorm:"column:course_id"`
	MoodleResource []MoodleResource `gorm:"foreignKey:MoodleCourseSectionId"`
}

type MoodleResource struct {
	ID                  uint   `json:"id" gorm:"primarykey"`
	ResourceID          uint   `json:"resource_id" `
	Instance            uint   `json:"instance"`
	Description         string `json:"description"`
	URL                 string `json:"url"`
	Name                string `json:"name"`
	ModIcon             string `json:"modicon"`
	ModName             string `json:"modname"`
	ResourceContent     string `json:"resource_content"`
	AssignIntroResource string `json:"assign_intro_resource"`

	MoodleCourseSectionId uint                             `json:"moodlecoursesectionid" gorm:"column:section_id"`
	Accounts              []MoodleAccount                  `gorm:"many2many:moodle_account_resources;"`
	Submissions           []MoodleAssignSubmissionResource `gorm:"foreignKey:MoodleResourceID"`
}
type MoodleAssignSubmissionResource struct {
	ID                         uint   `json:"id" gorm:"primarykey"`
	AssignSubmissionResourceId uint   `json:"assign_submission_resource_id"`
	GradeForDisplay            string `json:"gradefordisplay"`
	GradedDate                 int    `json:"gradeddate"`
	FeedbackComment            string `json:"feedbackcomment"`
	SubmissionAttachments      string `json:"submissionattachments"`
	AccountId                  uint   `json:"account_id"`
	MoodleResourceID           uint   `json:"moodle_resource_id" gorm:"column:resource_id"`
}

type MoodleAssignIntroResource struct {
	Intro               string `json:"intro"`
	SubmissionStatement string `json:"submissionstatement"`
	IntroAttachment     string `json:"introattachments"`
}
