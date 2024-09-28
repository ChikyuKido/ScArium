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
	ComponentID int             `json:"id"`
	ID          int             `json:"instance"`
	Description string          `json:"description"`
	URL         string          `json:"url"`
	Name        string          `json:"name"`
	ModIcon     string          `json:"modicon"`
	ModName     string          `json:"modname"`
	Contents    []CourseContent `json:"contents"`
}

type CourseContent struct {
	Type        string `json:"type"`
	FileName    string `json:"filename"`
	FileSize    int64  `json:"filesize"`
	FileURL     string `json:"fileurl"`
	TimeCreated int64  `json:"timecreated"`
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
