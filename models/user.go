package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type UserInfo struct {
	gorm.Model
	userName string
	Password string
	age      int
	sex      int
	birthday time.Time
	email    string
	phone    string
}
