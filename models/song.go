package models

type Song struct {
	ID      int    `json:"id" gorm:"primaryKey;autoIncrement"` // 主键
	Title   string `json:"title" gorm:"size:255;not null"`     // 歌曲标题
	File_id string `json:"file_id" gorm:"size:255;not null"`   // 文件标识符
}
