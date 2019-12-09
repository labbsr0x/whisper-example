package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	"github.com/labbsr0x/goh/gohtypes"
)

// writePage loads a page using templates
func writePage(w http.ResponseWriter, pageName string, page interface{}) {
	buf := new(bytes.Buffer)
	content := template.Must(template.ParseFiles("./assets/html/" + pageName + ".html"))

	err := content.Execute(buf, page)
	gohtypes.PanicIfError("Unable to load page", http.StatusInternalServerError, err)

	_, err = w.Write(buf.Bytes())
	gohtypes.PanicIfError("Unable to render", http.StatusInternalServerError, err)
}

// setCookie set a simple cookie that expires in a week
func setCookie(w http.ResponseWriter, name, value string) {
	oneWeekFromNow := time.Now().Add(7 * 24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: oneWeekFromNow,
	})
}

// unsetCookie overwrite a cookie and make it expired
func unsetCookie(w http.ResponseWriter, name string) {
	past := time.Unix(0, 0)

	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Value:   "",
		Expires: past,
	})
}
