package main


import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const DIR = "./images"

func main() {
	http.HandleFunc("/images", list_images)
	http.HandleFunc("/image/", get_image)
	http.HandleFunc("/remove/", remove_image)
	http.HandleFunc("/add", add_image)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

type Image struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Error interface {
	Error() string
}

func send_error(w http.ResponseWriter, err Error, log_message ...any) {
        log.Println(log_message...)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))

}

func list_images(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	files, err := ioutil.ReadDir(DIR)
	if err != nil {
		send_error(w, err, "Error: Could not read directory")
		return
	}

	images := []Image{}
	for _, file := range files {
		file_name := file.Name()
		if matched, _ := regexp.MatchString(".(jpg|jpeg)$", file_name); matched {
                        image_name := strings.Replace(strings.Split(file_name, ".")[0], "_", " ", -1)
			images = append(images, Image{Name: image_name, URL: "/image/" + file_name})
		}
	}

	res, err := json.Marshal(images)

	if err != nil {
		send_error(w, err, "Error: Could not marshal images")
		return
	}

	w.Write(res)
}

func get_image(w http.ResponseWriter, r *http.Request) {
        log.Println("Received request for image: " + r.URL.Path)

        file_name := DIR + "/" + strings.Replace(r.URL.Path, "/image/", "", 1)
        log.Println("Serving filename" + file_name)

	http.ServeFile(w, r, file_name)
}

func remove_image(w http.ResponseWriter, r *http.Request) {
        file_name := DIR + "/" + strings.Replace(r.URL.Path, "/remove/", "", 1) + ".jpg"
	err := os.Remove(file_name)
	if err != nil {
		send_error(w, err, "Error: Could not remove file", file_name)
		return
	}
        log.Println("Removed file: " + file_name)
}

func add_image(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		send_error(w, err, "Error: Could not parse form")
		return
	}

	title := r.FormValue("title")

	if title == "" {
		send_error(w, errors.New("Title is required"), "Error: No title for image")
		return
	}
	log.Println("Received image: " + title)

	// Get image from form
	image, _, err := r.FormFile("image")
	if err != nil {
		send_error(w, err, "Error: Could not get image from form")
		return
	}
	defer image.Close()

	// Create new file
        file_name := DIR + "/" + strings.Replace(title, " ", "_", -1) + ".jpg"
	dst, err := os.Create(file_name)

	if err != nil {
		send_error(w, err, "Error: Could not create file")
		return
	}
	defer dst.Close()

	// Copy image to new file
	if _, err := io.Copy(dst, image); err != nil {
		send_error(w, err, "Error: Could not copy contents to file")
		return
	}
        log.Println("Image saved as: " + file_name)

	w.WriteHeader(http.StatusOK)
}
