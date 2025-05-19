package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/demola234/authentication/docs/statik"
	pb "github.com/demola234/authentication/infrastructure/api/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/status"
)

// Handle multipart form uploads
func HandleFileUpload(ctx context.Context, server pb.AuthServiceServer) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("Received upload request: %s %s", r.Method, r.URL.Path)

		// Parse multipart form
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Get userId
		userId := r.FormValue("userId")
		if userId == "" {
			log.Printf("Missing userId")
			http.Error(w, "userId is required", http.StatusBadRequest)
			return
		}

		// Get file
		file, fileHeader, err := r.FormFile("content")
		if err != nil {
			log.Printf("Error getting file: %v", err)
			http.Error(w, "Failed to get file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		log.Printf("Received file: %s, size: %d bytes", fileHeader.Filename, fileHeader.Size)

		// Read the file
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Call the gRPC handler
		resp, err := server.UploadImage(ctx, &pb.UploadImageRequest{
			UserId:  userId,
			Content: fileBytes,
		})

		if err != nil {
			log.Printf("Error from gRPC: %v", err)
			st, ok := status.FromError(err)
			if ok {
				http.Error(w, st.Message(), runtime.HTTPStatusFromCode(st.Code()))
			} else {
				http.Error(w, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// Serve a test upload page
func ServeTestUploadPage(w http.ResponseWriter, r *http.Request) {
	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Test Image Upload</title>
    </head>
    <body>
        <h1>Test Image Upload</h1>
        <form action="/api/v1/upload-image" method="POST" enctype="multipart/form-data">
            <div>
                <label for="userId">User ID:</label>
                <input type="text" id="userId" name="userId" required>
            </div>
            <div>
                <label for="content">Select Image:</label>
                <input type="file" id="content" name="content" required>
            </div>
            <button type="submit">Upload</button>
        </form>
    </body>
    </html>
    `
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// Debug upload handler
func HandleDebugUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("Received debug upload request")

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Printf("Error parsing form: %v", err)
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Form values: %v", r.MultipartForm.Value)
		log.Printf("File headers: %v", r.MultipartForm.File)

		// Get the file
		file, header, err := r.FormFile("content")
		if err != nil {
			log.Printf("Error getting file: %v", err)
			http.Error(w, "Failed to get file: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read a sample of the file
		buffer := make([]byte, 50)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("Error reading file: %v", err)
			http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"message":      "Debug information",
			"filename":     header.Filename,
			"size":         header.Size,
			"content_type": header.Header.Get("Content-Type"),
			"form_values":  r.MultipartForm.Value,
			"file_sample":  fmt.Sprintf("%v", buffer[:n]),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		html := `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Debug Upload</title>
        </head>
        <body>
            <h1>Debug File Upload</h1>
            <form action="/debug-upload" method="POST" enctype="multipart/form-data">
                <div>
                    <label for="userId">User ID:</label>
                    <input type="text" id="userId" name="userId" value="test123" required>
                </div>
                <div>
                    <label for="content">File:</label>
                    <input type="file" id="content" name="content" required>
                </div>
                <button type="submit">Test Upload</button>
            </form>
        </body>
        </html>
        `
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}
}
