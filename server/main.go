package main

import (
	"encoding/json"
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
	http.HandleFunc("/image/", image_handler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

type Image struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Error interface {
	Error() string
}

func send_error(w http.ResponseWriter, err Error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error()))

}

func list_images(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	files, err := ioutil.ReadDir(DIR)
	if err != nil {
		send_error(w, err)
		return
	}

	var images []Image
	for _, file := range files {
		file_name := file.Name()
		if matched, _ := regexp.MatchString(".(jpg|jpeg)$", file_name); matched {
			images = append(images, Image{Name: strings.Split(file_name, ".")[0], URL: "/image/" + file_name})
		}
	}
	res, err := json.Marshal(images)

	if err != nil {
		send_error(w, err)
		return
	}

	w.Write(res)
}

func image_handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, DIR+"/"+strings.Replace(r.URL.Path, "/image/", "", 1))
	case "DELETE":
                err := os.Remove(DIR + "/uploads/" + strings.Replace(r.URL.Path, "/image/", "", 1))
                if err != nil {
                        send_error(w, err)
                        return
                }
	case "POST":
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			send_error(w, err)
			return
		}

		title := r.FormValue("title")

		// Get image from form
		image, _, err := r.FormFile("image")
		if err != nil {
			send_error(w, err)
			return
		}
		defer image.Close()

		// Create new file
		dst, err := os.Create(
			DIR + "/uploads/" + strings.Replace(
				strings.ToLower(title), " ", "_", -1,
			) + ".jpg",
		)

		if err != nil {
			send_error(w, err)
			return
		}
		defer dst.Close()

		// Copy image to new file
		if _, err := io.Copy(dst, image); err != nil {
			send_error(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
