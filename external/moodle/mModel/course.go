package mModel

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
	ComponentID int                `json:"id"`
	ID          int                `json:"instance"`
	Description string             `json:"description"`
	URL         string             `json:"url"`
	Name        string             `json:"name"`
	ModIcon     string             `json:"modicon"`
	ModName     string             `json:"modname"`
	Dates       []CourseModuleDate `json:"dates"`
	Contents    []CourseContent    `json:"contents"`
}

type CourseContent struct {
	Type     string `json:"type"`
	FileName string `json:"filename"`
	FileSize int64  `json:"filesize"`
	FileURL  string `json:"fileurl"`
}

type CourseModuleDate struct {
	Label     string `json:"label"`
	Timestamp int64  `json:"timestamp"`
	DataID    string `json:"dataid"`
}
