package main

import (
	"fmt"
	"log"
	"mime"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"strings"
)

type EmptyArgs struct{}

type FileListResponse struct {
	Files   []string
	Success bool
	Message string
}

type FileRequest struct {
	Filename string
}

type FileResponse struct {
	Filename string
	Content  []byte
	MIMEType string
	Success  bool
	Message  string
}

type FileService struct {
	FilesDir  string
	ImagesDir string
}

func (f *FileService) ListFiles(args EmptyArgs, res *FileListResponse) error {
	var allFiles []string

	fileEntries, err := os.ReadDir(f.FilesDir)
	if err == nil {
		for _, entry := range fileEntries {
			if !entry.IsDir() {
				allFiles = append(allFiles, "files/"+entry.Name())
			}
		}
	}

	imageEntries, err := os.ReadDir(f.ImagesDir)
	if err == nil {
		for _, entry := range imageEntries {
			if !entry.IsDir() {
				allFiles = append(allFiles, "images/"+entry.Name())
			}
		}
	}

	if len(allFiles) == 0 {
		res.Success = false
		res.Message = "no files or images found"
		return nil
	}

	res.Files = allFiles
	res.Success = true
	res.Message = "files listed successfully"
	return nil
}

func (f *FileService) GetFile(req FileRequest, res *FileResponse) error {
	if req.Filename == "" {
		res.Success = false
		res.Message = "filename is required"
		return nil
	}

	cleanPath := filepath.Clean(req.Filename)

	var fullPath string

	if strings.HasPrefix(cleanPath, "files/") {
		name := filepath.Base(cleanPath)
		fullPath = filepath.Join(f.FilesDir, name)
	} else if strings.HasPrefix(cleanPath, "images/") {
		name := filepath.Base(cleanPath)
		fullPath = filepath.Join(f.ImagesDir, name)
	} else {
		res.Success = false
		res.Message = "invalid file path"
		return nil
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		res.Success = false
		res.Message = "file not found or cannot be read"
		return nil
	}

	ext := filepath.Ext(fullPath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	res.Filename = req.Filename
	res.Content = data
	res.MIMEType = mimeType
	res.Success = true
	res.Message = "file retrieved successfully"

	log.Printf("File sent via JSON-RPC: %s", req.Filename)
	return nil
}

func main() {
	fileService := &FileService{
		FilesDir:  "files",
		ImagesDir: "images",
	}

	err := rpc.Register(fileService)
	if err != nil {
		log.Fatalf("Failed to register file RPC service: %v", err)
	}

	listener, err := net.Listen("tcp", ":2345")
	if err != nil {
		log.Fatalf("Failed to listen on port 2345: %v", err)
	}

	fmt.Println("File JSON-RPC server running on port 2345...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}
