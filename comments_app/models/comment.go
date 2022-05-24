package models

import "time"

type Comment struct {
	Id             uint      `json:"id" gorm:"primaryKey"`
	Content        string    `json:"content"`
	PostId         uint64    `json:"-"`
	Timestamp      time.Time `json:"timestamp"`
	AuthorId       uint64    `json:"authorId"`
	AuthorUsername string    `json:"authorUsername"`
}
