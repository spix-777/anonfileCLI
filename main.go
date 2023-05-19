package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// Loggers
var (
	BannerLogger *log.Logger
	InfoLogger   *log.Logger
	WarnLogger   *log.Logger
)

// Initialize loggers to stdout
func init() {
	BannerLogger = log.New(os.Stdout, "", 0)
	InfoLogger = log.New(os.Stdout, " [ + ] ", 0)
	WarnLogger = log.New(os.Stdout, " [ ! ] ", 0)

}

func main() {
	var fileStr string
	var infoStr string
	var version bool
	flag.StringVar(&fileStr, "f", "", "File to upload")
	flag.StringVar(&infoStr, "i", "", "Info about the file/ID")
	flag.BoolVar(&version, "v", false, "Version")
	flag.Parse()

	// Print the banner, of course!
	banner()

	// Version
	if version == true {
		InfoLogger.Println("Version: 1.0.0")
		os.Exit(0)
	}

	// Upload file or get info
	if fileStr != "" {
		InfoLogger.Println("Uploading file:", fileStr)
		fileUpload(fileStr)
	} else if infoStr != "" {
		InfoLogger.Println("Getting info about file:", infoStr)
		fileInfo(infoStr)
	}

}

func fileInfo(infoStr string) {
	infoURL := "https://api.anonfiles.com/v2/file/" + infoStr + "/info"
	resp, err := http.Get(infoURL)
	if err != nil {
		WarnLogger.Println("Error getting file info:", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		WarnLogger.Println("Error reading response:", err)
		return
	}
	//InfoLogger.Println(string(responseBody))
	if strings.Index(string(responseBody), "error") != -1 {
		WarnLogger.Println("Error getting file info:", infoStr)
		return
	}
	// Get the short URL
	reqStr := string(responseBody)
	indexShort := strings.Index(reqStr, "short")
	start := reqStr[indexShort+9:]
	indexShort = strings.Index(start, `",`)
	short := start[:indexShort]
	out := strings.Replace(short, "\"}", "", -1)
	out = strings.Replace(short, "\\", "", -1)
	InfoLogger.Println("Short URL:", out)

	// Get the size
	indexSize := strings.Index(reqStr, "bytes")
	start = reqStr[indexSize+7:]
	indexSize = strings.Index(start, `,`)
	size := start[:indexSize]
	InfoLogger.Println("Size:", size+" bytes")

	// Get the status
	indexStatus := strings.Index(reqStr, "status")
	start = reqStr[indexStatus+8:]
	indexStatus = strings.Index(start, `,"`)
	status := start[:indexStatus]
	InfoLogger.Println("Status:", status)

}

func fileUpload(fileStr string) {
	// Open the file to be uploaded
	file, err := os.Open(fileStr)
	if err != nil {
		WarnLogger.Println("Error opening file ...")
		return
	}
	defer file.Close()

	// Create a new multipart form writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a new form field for the file
	fileField, err := writer.CreateFormFile("file", fileStr)
	if err != nil {
		WarnLogger.Println("Error creating form field ...")
		return
	}

	// Copy the file contents to the form field
	_, err = io.Copy(fileField, file)
	if err != nil {
		WarnLogger.Println("Error copying file contents ...")
		return
	}

	// Close the multipart form writer
	err = writer.Close()
	if err != nil {
		WarnLogger.Println("Error closing form writer ...")
		return
	}

	// Create a new POST request with the form data
	url := "https://api.anonfiles.com/upload"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		WarnLogger.Println("Error creating request ...")
		return
	}

	// Set the content type header for the form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		WarnLogger.Println("Error sending request ...")
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		WarnLogger.Println("Error reading response ...")
		return
	}

	// Get the short URL
	reqStr := string(responseBody)
	reqStr = strings.ReplaceAll(reqStr, " ", "")
	indexShort := strings.Index(reqStr, "short")
	start := reqStr[indexShort+8:]
	indexShort = strings.Index(start, `"`)
	short := start[:indexShort]
	out := strings.Replace(short, "\"}", "", -1)
	out = strings.Replace(short, "\\", "", -1)
	InfoLogger.Println("Short URL:", out)

	// Get the size
	indexSize := strings.Index(reqStr, "readable")
	start = reqStr[indexSize+11:]
	indexSize = strings.Index(start, `B`)
	size := start[:indexSize]
	InfoLogger.Println("Size:", size+" bytes")

	// Get the id
	indexID := strings.Index(reqStr, "id")
	start = reqStr[indexID+6:]
	indexID = strings.Index(start, `",`)
	id := start[:indexID]
	InfoLogger.Println("ID:", id)

}

func banner() {
	red := "\033[31m"
	reset := "\033[0m"
	love := red + "<3" + reset
	banner := `                                             
                          __ _ _           
  __ _ _ __   ___  _ __  / _(_) | ___  ___ 
 / _  |  _ \ / _ \| '_ \| |_| | |/ _ \/ __|
| (_| | | | | (_) | | | |  _| | |  __/\__ \
 \__,_|_| |_|\___/|_| |_|_| |_|_|\___||___/
	`
	BannerLogger.Print(banner)
	BannerLogger.Println(strings.Repeat(" ", 23) + "made by " + love + " @SpiX-777")
}
