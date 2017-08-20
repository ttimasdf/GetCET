package handler

import (
	"html/template"
	"net/http"
)

// Score 类型定义了成绩查询返回值所包含的字段
type Score struct {
	Name, School, Listening, Reading, Writing, Total, Error string
}

var tmpl *template.Template

// ScoreHandler 用于查询成绩
func ScoreHandler(w http.ResponseWriter, req *http.Request) {
	if tmpl == nil {
		tmpl = template.Must(template.ParseFiles("public/result.html"))
	}
	tmpl.Execute(w, Score{"haha", "haha", "haha", "haha", "haha", "haha", ""})
}
