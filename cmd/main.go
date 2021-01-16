package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(20 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)

	fmt.Printf("File Size: %+v\n", handler.Size)

	fmt.Printf("MIME Header: %+v\n", handler.Header)
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.tgs")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error read file content")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// write this byte array to our temporary file
	_, err = tempFile.Write(fileBytes)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error save the File")
		fmt.Println(err)
		return
	}
	tempFile.Close()
	name := tempFile.Name()
	convertedName := name + ".webp"
	cmd := exec.Command("python3", "convert.py", name, convertedName)

	defer func() {
		_ = os.Remove(name)
		_ = os.Remove(convertedName)
	}()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, convertedName)
}
func main() {

	http.HandleFunc("/upload", uploadFile)
	fmt.Println("starting server")
	http.ListenAndServe(":8080", nil)
}
