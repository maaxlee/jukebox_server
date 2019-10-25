package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const credFile = "creds.txt"

var (
	songsNum     = 10
	songsListMap = make(map[string][]string)
	usersMap     = make(map[string]string)
)

func getSong(user string) string {

	log.Println(songsNum)
	songsList := songsListMap[user]
	if len(songsList) == 0 {
		log.Print("List is empty, create a new list")
		for i := 1; i < songsNum+1; i++ {
			songsList = append(songsList, "Song number "+strconv.Itoa(i))
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(songsList), func(i, j int) { songsList[i], songsList[j] = songsList[j], songsList[i] })
	}

	ret := songsList[0]
	songsList = songsList[1:]
	songsListMap[user] = songsList
	log.Println("User: " + user)
	log.Println(songsList)
	return ret
}

func checkCreds(u, p string) bool {
	if pass, ok := usersMap[u]; !ok || p != pass {
		return false
	}
	return true

}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || !checkCreds(user, pass) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Unauthorized!")
		return
	}
	fmt.Fprintln(w, getSong(user))
}

func main() {

	// Get songs number
	if n, ok := os.LookupEnv("SONGS_NUM"); ok {
		var err error
		songsNum, err = strconv.Atoi(n)
		if err != nil {
			log.Fatal("Wrong songs number format type")
		}
	}

	// Reading credentials file
	log.Println("Reading DB")
	data, err := ioutil.ReadFile(credFile)
	if err != nil {
		log.Fatal("Could not read credentials DB")
	}
	strList := strings.Split(string(data), "\n")
	for _, str := range strList {
		if str != "" {
			raw := strings.Split(str, ":")
			usersMap[raw[0]] = raw[1]
			songsListMap[raw[0]] = make([]string, 0)
		}
	}

	// start server
	log.Print("Starting server")
	http.HandleFunc("/", simpleHandler)
	http.ListenAndServe(":8080", nil)
}
