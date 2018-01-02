package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	linkdbPath string = "./linkdb"
)

// link is a dumb type representing a shortened link
type link string

func (l link) String() string {
	return string(l)
}

// linkdb is a small structure which holds all the shortlinks,
// handles on-disk retention.
type linkdb struct {
	// memoized, shorted Urls
	Urls map[string]link `json:"Urls"`
}

// LoadLinkdbFromDisk loads a linkdb from the standard path
func LoadLinkdbFromDisk() (linkdb, error) {
	data, err := ioutil.ReadFile(linkdbPath)
	if err != nil {
		// just initialize a new linkdb and return it
		ld := linkdb{}
		ld.Urls = make(map[string]link)

		return ld, nil
	}

	var ld linkdb
	err = json.Unmarshal(data, &ld)

	// if linkdb is empty even after the unmarshaling, initialize inner struct
	if len(ld.Urls) == 0 {
		ld.Urls = make(map[string]link)
	}

	return ld, err
}

// saveOnDisk saves ld on the linkdb file
func (ld linkdb) saveOnDisk() error {
	data, err := json.Marshal(ld)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(linkdbPath, data, os.FileMode(0644))
	return err
}

// shortAndRetain checks if url is a valid net.URL, shorts and memorizes if valid
func (ld *linkdb) shortAndRetain(urlStr string) (string, error) {
	// parse urlStr
	_, err := url.Parse(urlStr)

	if err != nil {
		// invalid url, raise error
		return "", err
	}

	// url is valid, hash it and memorize
	hash := shortHash(urlStr)
	ld.Urls[hash] = link(shallIPrependHTTP(urlStr))

	// initiate linkdb on-disk retain
	err = ld.saveOnDisk()
	if err != nil {
		log.Println(err)
	}

	return hash, nil
}

// get gets a link by its hash, or gives error
func (ld linkdb) get(hash string) (link, error) {
	if u, ok := ld.Urls[hash]; ok {
		return u, nil
	}

	return link(""), errors.New("shortlink not found")
}

// shallIPrependHTTP blindly adds "http://" if the string doesn't have any
func shallIPrependHTTP(s string) string {
	if !(strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")) {
		return "http://" + s
	}

	return s
}

// shortHash returns 5 sequential chars from the hashed string,
// starting from the 5th byte
func shortHash(s string) string {
	h := sha256.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		log.Println("hashing error!", err)
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return string(sum[5:10])
}
