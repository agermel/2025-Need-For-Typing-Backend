package models

type Asset struct {
	ID      int    `json:"id" `      // 主键
	Name    string `json:"name" `    // 资源名称
	File_id string `json:"file_id" ` // 文件标识符
}
