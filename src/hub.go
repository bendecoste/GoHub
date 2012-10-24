package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	API_ISSUES = "https://api.github.com/issues"
	API_ISSUE  = "https://api.github.com/repos/goinstant/slashjoin/issues/"
	ACCESS     = "?access_token=10520ec3d313a78e387fe9cf54da6f2af9acdbf4"
)

func comment(num string) {
	res, err := http.Get(API_ISSUE + num + "/comments" + ACCESS)

	if err != nil {
		fmt.Println("Fatal Error on ISSUE", err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	type Comment struct {
		Body string
		User struct {
			Login string
		}
	}

	var c []Comment
	err = json.Unmarshal(resp, &c)

	if err != nil {
		fmt.Println("Fatal Comment --", err)
	}

	fmt.Printf(" -- Comments\n")
	for _, val := range c {
		fmt.Printf("\tUser: %s\n", val.User.Login)
		fmt.Printf("\tBody: %s\n\n", val.Body)
	}
}

func issue(num string) {
	res, err := http.Get(API_ISSUE + num + ACCESS)

	if err != nil {
		fmt.Println("Fatal Error on ISSUE", err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	type Message struct {
		Html_url string
		Title    string
		Body     string
	}

	var m Message
	err = json.Unmarshal(resp, &m)
	if err != nil {
		fmt.Println("Fatal Error", err)
	}

	fmt.Printf(" -- Title\n\t%s\n\n -- Body\n\t%s\n\n -- URL\n\t%s\n", m.Title, m.Body, m.Html_url)
	comment(num)
}

func issues() {
	res, err := http.Get(API_ISSUES + ACCESS)

	if err != nil {
		fmt.Println("Fatal Error on ISSUES", err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	type Message struct {
		Url      string
		Number   int
		Comments int
		Title    string
	}

	// do something with error?
	var m []Message
	err = json.Unmarshal(resp, &m)
	if err != nil {
		fmt.Println("Fatal Error on JSON.UNMARSHAL")
	}

	for _, val := range m {
		fmt.Printf(" -- #%d > %s\n", val.Number, val.Title)
	}
}

func execute(cmd, num string) {
	switch cmd {
	case "issues":
		issues()
		return
	case "issue":
		issue(num)
		return
    case "comments":
        comment(num)
        return
	}
}

func main() {
	var cmd = flag.String("cmd", "", "Help")
	var issueNum = flag.String("num", "", "Help")
	flag.Parse()

	execute(*cmd, *issueNum)
}
