package render

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var (
	path      string
	templates map[string]*template.Template
)

func SetPath(p string) {
	path = strings.TrimRight(p, "/") + "/"
}

func JSON(w http.ResponseWriter, data interface{}, code ...int) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// set HTTP response code
	c := http.StatusOK
	if len(code) > 0 {
		c = code[0]
	}

	w.WriteHeader(c)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return nil
}

func File(w http.ResponseWriter, file string, context map[string]interface{}, code ...int) error {
	return Files(w, []string{file}, context, code...)
}

func Files(w http.ResponseWriter, files []string, context map[string]interface{}, code ...int) error {
	key := strings.Join(files, "\n")
	templateFiles := make([]string, len(files))
	for i, v := range files {
		templateFiles[i] = fmt.Sprintf("%s%s", path, v)
	}

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

	w.WriteHeader(c)
	w.Header().Set("Content-Type", "text/html")
	return t.Execute(w, context)
}

func FileInLayout(w http.ResponseWriter, layout, file string, context map[string]interface{}, code ...int) error {
	return Files(w, []string{layout, file}, context, code...)
}

func FilesInLayout(w http.ResponseWriter, layout string, files []string, context map[string]interface{}, code ...int) error {
	f := []string{layout}
	f = append(f, files...)
	return Files(w, f, context, code...)
}

func Error(w http.ResponseWriter, file string, code int) {
	if err := Files(w, []string{file}, map[string]interface{}{}, code); err != nil {
		http.Error(w, http.StatusText(code), code)
		return
	}
}