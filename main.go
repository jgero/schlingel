package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/upload", uploadHandler(minioClient))
	prx := createS3ReverseProxy(minioClient)
	http.Handle("/download/", prx)
	// http.HandleFunc("/download/", downloadHandler(minioClient))

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func uploadHandler(client *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit upload size
		r.ParseMultipartForm(32 << 20) // 32MB

		// Get file from form
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get file from request", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Generate a key for the file (you might want a better strategy)
		key := fmt.Sprintf("%d-%s", time.Now().Unix(), header.Filename)

		// Upload to S3
		_, err = client.PutObject(r.Context(), "bucket", key, file, -1, minio.PutObjectOptions{})
		if err != nil {
			log.Printf("Failed to upload file: %v", err)
			http.Error(w, "Failed to store file", http.StatusInternalServerError)
			return
		}

		// Return the key or path to access the file
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"success": true, "file_path": "/download/%s"}`, key)
	}
}

func downloadHandler(client *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract file key from request, e.g., from URL path or query parameter
		key := r.URL.Path[len("/download/"):]
		if key == "" {
			http.Error(w, "File key is required", http.StatusBadRequest)
			return
		}

		stat, err := client.StatObject(r.Context(), "bucket", key, minio.StatObjectOptions{})
		if err != nil {
			log.Printf("Failed to stat object: %v", err)
			http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
			return
		}
		obj, err := client.GetObject(r.Context(), "bucket", key, minio.GetObjectOptions{})
		if err != nil {
			log.Printf("Failed to get object: %v", err)
			http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
			return
		}
		defer obj.Close()

		// Set appropriate headers
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", key))
		w.Header().Set("Content-Type", stat.ContentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size))

		// Stream the file directly to the client
		_, err = io.Copy(w, obj)
		if err != nil {
			log.Printf("Failed to stream file: %v", err)
		}
	}
}

func createS3ReverseProxy(client *minio.Client) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		// Extract key from the original request path
		key := req.URL.Path[len("/download/"):]

		presignedUrl, err := client.PresignedGetObject(req.Context(), "bucket", key, 15*time.Minute, url.Values{})
		if err != nil {
			log.Printf("Failed to generate presigned URL: %v", err)
			return
		}
		req.URL.Scheme = presignedUrl.Scheme
		req.URL.Host = presignedUrl.Host
		req.URL.Path = presignedUrl.Path
		req.URL.RawQuery = presignedUrl.RawQuery
		req.Host = presignedUrl.Host
	}

	return &httputil.ReverseProxy{
		Director: director,
	}
}
