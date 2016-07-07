package render

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	path      string
	templates map[string]*template.Template
)

func SetPath(p string) {
	path = strings.TrimRight(p, string(os.PathSeparator)) + string(os.PathSeparator)
}

func JSON(w http.ResponseWriter, data interface{}, code ...int) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Panic(err)
		return
	}

	// set HTTP response code
	c := http.StatusOK
	if len(code) > 0 {
		c = code[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	w.Write(b)
}

func File(w http.ResponseWriter, file string, context map[string]interface{}, code ...int) {
	if err := renderFiles(w, []string{file}, context, code...); err != nil {
		log.Panic(err)
	}
}

func Files(w http.ResponseWriter, files []string, context map[string]interface{}, code ...int) {
	if err := renderFiles(w, files, context, code...); err != nil {
		log.Panic(err)
	}
}

func FileInLayout(w http.ResponseWriter, layout, file string, context map[string]interface{}, code ...int) {
	if err := renderFiles(w, []string{file, layout}, context, code...); err != nil {
		log.Panic(err)
	}
}

func FilesInLayout(w http.ResponseWriter, layout string, files []string, context map[string]interface{}, code ...int) {
	files = append(files, layout)
	if err := renderFiles(w, files, context, code...); err != nil {
		log.Panic(err)
	}
}

func Error(w http.ResponseWriter, file string, code int) {
	if err := renderFiles(w, []string{file}, map[string]interface{}{}, code); err != nil {
		http.Error(w, http.StatusText(code), code)
	}
}

func renderFiles(w http.ResponseWriter, files []string, context map[string]interface{}, code ...int) error {
	templateFiles := make([]string, len(files))
	for i, v := range files {
		templateFiles[i] = fmt.Sprintf("%s%s", path, v)
	}
	key := strings.Join(templateFiles, "\n")

	var t *template.Template
	var ok bool
	if t, ok = templates[key]; !ok {
		var err error
		t, err = template.ParseFiles(templateFiles...)
		if err != nil {
			return err
		}
	}

	// set HTTP response code
	c := http.StatusOK
	if len(code) > 0 {
		c = code[0]
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(c)
	return t.Execute(w, context)
}
