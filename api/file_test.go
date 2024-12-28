package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	t.Run("upload without attached file", func(t *testing.T) {
		router := BuildRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/files/upload", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "", w.Body.String())
	})
}

func TestDownloadFile(t *testing.T) {
	t.Run("download non-existent file", func(t *testing.T) {
		router := BuildRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/files/foo.pdf", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "'foo.pdf' does not exist.", w.Body.String())
	})

	t.Run("download existent file", func(t *testing.T) {
		router := BuildRouter()

		filename := "testfile.pdf"
		file, err := os.Create(filename)
		assert.NoError(t, err)
		_, err = file.WriteString("TEST FILE CONTENT TROLOLOL")
		assert.NoError(t, err)
		file.Close()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/files/%s", filename), nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "TEST FILE CONTENT TROLOLOL'testfile.pdf' downloaded!", w.Body.String())

		err = os.Remove(filename)
		assert.NoError(t, err)
	})
}
