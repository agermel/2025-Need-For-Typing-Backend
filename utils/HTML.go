package utils

import (
	"io/ioutil"
	"net/http"
)

// serveHTML 提供静态HTML文件
func ServeHTML(w http.ResponseWriter, r *http.Request, route string) {
	data, err := ioutil.ReadFile(route)
	if err != nil {
		http.Error(w, "文件读取错误", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
