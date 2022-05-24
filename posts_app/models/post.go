package models

import "time"

type Post struct {
	Id             uint      `json:"id" gorm:"primaryKey"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
	AuthorId       uint64    `json:"authorId"`
	AuthorUsername string    `json:"authorUsername"`
}
