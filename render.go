package render

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

var (
	path      string
	templates map[string]*template.Template

	funcMap    = template.FuncMap{}
	rxFuncName = regexp.MustCompile("^[a-zA-Z]+$")
)

var (
	ErrNoFile          = errors.New("render: no files named in call to execute")
	ErrInvalidFuncName = errors.New("render: function name can only have alpha characters")
	ErrNotFunc         = errors.New("render: function required in call to SetFunc")
)

func SetPath(p string) {
	path = strings.TrimRight(p, string(os.PathSeparator)) + string(os.PathSeparator)
}

func SetFunc(name string, f interface{}) error {
	if reflect.ValueOf(f).Kind() != reflect.Func {
		return ErrNotFunc
	}
	if !rxFuncName.MatchString(name) {
		return ErrInvalidFuncName
	}
	funcMap[name] = f
	return nil
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
	if err := execute(w, []string{file}, context, code...); err != nil {
		log.Panic(err)
	}
}

func Files(w http.ResponseWriter, files []string, context map[string]interface{}, code ...int) {
	if err := execute(w, files, context, code...); err != nil {
		log.Panic(err)
	}
}

func FileInLayout(w http.ResponseWriter, layout, file string, context map[string]interface{}, code ...int) {
	if err := execute(w, []string{file, layout}, context, code...); err != nil {
		log.Panic(err)
	}
}

func FilesInLayout(w http.ResponseWriter, layout string, files []string, context map[string]interface{}, code ...int) {
	files = append(files, layout)
	if err := execute(w, files, context, code...); err != nil {
		log.Panic(err)
	}
}

func Error(w http.ResponseWriter, file string, code int) {
	if err := execute(w, []string{file}, map[string]interface{}{}, code); err != nil {
		http.Error(w, http.StatusText(code), code)
	}
}

func execute(w http.ResponseWriter, files []string, context map[string]interface{}, code ...int) error {
	if len(files) == 0 {
		return ErrNoFile
	}

	templateFiles := make([]string, len(files))
	for i, v := range files {
		templateFiles[i] = fmt.Sprintf("%s%s", path, v)
	}
	key := strings.Join(templateFiles, "\n")

	var t *template.Template
	var ok bool
	if t, ok = templates[key]; !ok {
		var err error
		t, err = template.New(filepath.Base(templateFiles[0])).
			Funcs(funcMap).
			ParseFiles(templateFiles...)
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
