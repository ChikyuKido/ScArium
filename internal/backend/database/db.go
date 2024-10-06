package database

import (
	"ScArium/common/config"
	"ScArium/common/log"
	"ScArium/internal/backend/database/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(config.DATABASE_PATH), &gorm.Config{})
	if err != nil {
		log.I.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		log.I.Fatalf("failed to create user database: %v", err)
	}
	err = db.AutoMigrate(&entity.D4sAccount{})
	if err != nil {
		log.I.Fatalf("failed to create d4s account database: %v", err)
	}
	err = db.AutoMigrate(&entity.MoodleAccount{}, &entity.MoodleCourse{}, &entity.MoodleCourseSection{}, &entity.MoodleResource{}, &entity.MoodleAssignSubmissionResource{})
	if err != nil {
		log.I.Fatalf("failed to create moodle account database: %v", err)
	}
}

func GetDB() *gorm.DB {
	return db
}
