package controller

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

func StartPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//указываем путь к нужному файлу
	path := filepath.Join("public", "html", "index.html")
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
	//выводим шаблон клиенту в браузер
	err = tmpl.Execute(rw, nil)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}
}
