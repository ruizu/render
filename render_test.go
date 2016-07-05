package render

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		File(w, "sample/content.single.html", map[string]interface{}{"hello": "Ruizu"}, http.StatusOK)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal("http test get error: ", err)
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("http test read data error: ", err)
	}

	expected := `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Render File</title>
</head>

<body>
File content
</body>
</html>
`
	if string(data) != expected {
		t.Fatalf("http test unexpected data: %q.", data)
	}
}

func TestRenderError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Error(w, "sample/content.single.html", http.StatusInternalServerError)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal("http test get error: ", err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatal("http test invalid status response: ", http.StatusText(res.StatusCode))
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("http test read data error: ", err)
	}

	expected := `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Render File</title>
</head>

<body>
File content
</body>
</html>
`
	if string(data) != expected {
		t.Fatalf("http test unexpected data: %q.", data)
	}
}

func TestRenderErrorNoFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Error(w, "sample/error-no-file.html", http.StatusInternalServerError)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal("http test get error: ", err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatal("http test invalid status response: ", http.StatusText(res.StatusCode))
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal("http test read data error: ", err)
	}

	expected := "Internal Server Error\n"
	if string(data) != expected {
		t.Fatalf("http test unexpected data: %q.", data)
	}
}
