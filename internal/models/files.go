package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	UID  string `gorm:"unique;index" json:"uid"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type Chunks struct {
	gorm.Model
	File       File   `json:"file" gorm:"foreignKey:UID"`
	Name       string `json:"name"`
	FileUID    string `json:"file_uid"`
	DiscordURL string `json:"discord_url"`
	Order      int    `json:"order"`
}
