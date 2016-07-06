package render

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetPath("sample")
	os.Exit(m.Run())
}

func TestRenderFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		File(w, "content.single.html", map[string]interface{}{}, http.StatusOK)
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

func TestRenderFileInLayout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		FileInLayout(w, "layout.html", "content.html", map[string]interface{}{}, http.StatusOK)
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

	expected := `
<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Render File in Layout</title>
</head>

<body>

Content inside layout

</body>
</html>





`
	if string(data) != expected {
		t.Fatalf("http test unexpected data: %q.", data)
	}
}

func TestRenderError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Error(w, "content.single.html", http.StatusInternalServerError)
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
		Error(w, "error-no-file.html", http.StatusInternalServerError)
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
