package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

type BusyTime struct {
	Start time.Time
	End   time.Time
}

func fetchEvents(client *http.Client) (times []*BusyTime) {
	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)

	fmt.Println(t)
	max := time.Now().Add(time.Hour * 2).Format(time.RFC3339)

	// Fetch the user's calendars. TODO: allow user to specify which calendars to monitor for busy-ness.
	ls, err := srv.CalendarList.List().ShowDeleted(false).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve users calendars %v", err)
	}

	var calIds []*calendar.FreeBusyRequestItem
	for _, item := range ls.Items {
		// fmt.Println(item.Id)
		calIds = append(calIds, &calendar.FreeBusyRequestItem{Id: item.Id})
	}

	res, err := srv.Freebusy.Query(&calendar.FreeBusyRequest{TimeMin: t, TimeMax: max, Items: calIds}).Do()
	if err != nil {
		log.Fatalf("Unable to fetch free-busy user's events: %v", err)
	}

	fmt.Println()
	for _, cal := range res.Calendars {
		for _, busyTime := range cal.Busy {
			fmt.Println(busyTime.Start, busyTime.End)
			start, err := time.Parse(time.RFC3339, busyTime.Start)
			if err != nil {
				log.Println("Warn: failed to parse start time " + busyTime.Start)
			}
			end, err := time.Parse(time.RFC3339, busyTime.End)
			if err != nil {
				log.Println("Warn: failed to parse end time" + busyTime.End)
			}
			times = append(times, &BusyTime{start, end})
		}
	}
	fmt.Println()
	return
}

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	busyTimes := fetchEvents(client)
	fmt.Printf("%v events found\n", len(busyTimes))
	shouldBeOn := false
	now := time.Now()
	for _, bt := range busyTimes {
		if bt.Start.Before(now) && bt.End.After(now) {
			shouldBeOn = true
			break
		}
	}
	fmt.Printf("Should light be on: %v\n", shouldBeOn)
}
