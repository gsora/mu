package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

var (
	baseURL     string
	completeURL string
	port        string
	addPort     bool
	ld          linkdb
	err         error
)

func main() {
	setupParams()

	ld, err = LoadLinkdbFromDisk()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("running with baseurl", baseURL, "on port", port)

	// assemble baseURL
	if addPort {
		completeURL = "http://" + baseURL + ":" + port + "/"
	} else {
		completeURL = "http://" + baseURL + "/"
	}

	mux := http.NewServeMux()

	th := http.HandlerFunc(handleRequestByType)
	mux.Handle("/", th)

	log.Fatal(http.ListenAndServe(baseURL+":"+port, mux))
}

func setupParams() {
	flag.StringVar(&baseURL, "domain", "localhost", "domain to prepend to shortlinks")
	flag.StringVar(&port, "port", "8080", "port mu will bind on")
	flag.BoolVar(&addPort, "appendport", false, "append port to the link url")
	flag.Parse()
}

func handleRequestByType(w http.ResponseWriter, r *http.Request) {
	splitPath := strings.SplitN(r.URL.String(), "/", 3)

	// get the verb
	verb := strings.SplitN(splitPath[1], "?", 2)[0]

	// if user wants to add a new url, do that
	if verb == "add" {
		urlToShort := r.URL.Query().Get("url")

		if urlToShort == "" {
			httpWrite("missing link", w)
			return
		}

		hash, err := ld.shortAndRetain(urlToShort)
		if err != nil {
			httpWrite("invalid url", w)
			return
		}

		httpWrite(completeURL+hash, w)
		return
	}

	// otherwise, just assume shorturl fetching
	hash := splitPath[1]
	log.Println("fetching for hash", hash)
	url, err := ld.get(hash)
	log.Println("got url:", url.String())
	if err != nil {
		httpWrite("I don't know this hash!", w)
		return
	}

	// redirect to url
	http.Redirect(w, r, url.String(), 301)
}

// writes to w, and logs errors if any
func httpWrite(s string, w http.ResponseWriter) {
	_, err := w.Write([]byte(s))
	if err != nil {
		log.Println(err)
	}
}
