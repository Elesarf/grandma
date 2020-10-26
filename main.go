package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var allowedExtensionList = []string{"png", "jpg", "jpeg"}

var systemConfig ConfigData

func contains(arr []string, val string) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("pic")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileName := handler.Filename
	tmp := strings.Split(fileName, ".")
	extension := tmp[len(tmp)-1]

	if !contains(allowedExtensionList, extension) {
		fmt.Println("Error: not allowed extension")
		return
	}

	if !CheckPictureName(fileName) {
		fmt.Println("Error: not allowed file name")
		return
	}

	picVal, err := ConstructName(tmp[0])
	if err == nil {
		log.Println("Good name: " + picVal.UTCTime.String() + " : " + picVal.mainLocation + " : " + picVal.additionalLocation)
	} else {
		log.Println("Error parse name: " + err.Error())
	}

	log.Println("Uploaded File: " + fileName + ", size: " + strconv.FormatInt(handler.Size, 10) + " kB")

	resultFile, err := os.Create(systemConfig.pictureSavePath + "/" + fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer resultFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	resultFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	systemConfig = loadConfigFromEnv()
	fmt.Println("Start server")
	setupRoutes()
}
