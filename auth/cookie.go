package auth

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

var cookieStore = &CookieStore{}
var wd, _ = os.Getwd()
var cookiePath = path.Join(wd, "cookie")

type CookieStore struct {
	cookie string
}

func init() {
	if _, err := os.Stat(cookiePath); os.IsNotExist(err) {
		log.Println("CookieStore:", "Cookie file not found.")
		panic(err)
	}
	cookieStore.Reload()
}

func GetStore() *CookieStore {
	return cookieStore
}

func (*CookieStore) Reload() error {
	payload, err := ioutil.ReadFile(cookiePath)
	if err != nil {
		log.Println("CookieStore:", "Cannot read cookie file at", cookiePath)
		return err
	}
	cookieStore.cookie = string(payload)
	log.Printf("CookieStore: Reloaded to \n%s\n", cookieStore.cookie)
	return nil
}

func (s *CookieStore) Start() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("CookieStore:", "Fail to create file watcher for path", cookiePath, err)
	}

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("CookieStore: EVENT! %s\n", event.String())
				cookieStore.Reload()
				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("CookieStore: file watcher error: ERROR", err)
			}
		}
	}()

	if err := watcher.Add(cookiePath); err != nil {
		fmt.Println("CookieStore:", "Fail to add file to watcher cookiePath=", cookiePath, err)
	}
}
func (s *CookieStore) Cookie() string {
	return s.cookie
}
func (s *CookieStore) SetCookie(r *http.Request) {
	r.Header.Add("Cookie", s.cookie)
}
