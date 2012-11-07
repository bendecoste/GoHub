package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	API_URL           = "https://api.github.com"
	ACCESS            = "?access_token=10520ec3d313a78e387fe9cf54da6f2af9acdbf4"
	API_ISSUES        = "https://api.github.com/issues"
	API_ISSUE         = "https://api.github.com/repos/goinstant/slashjoin/issues/"
	API_COMMENT       = "https://api.github.com/repos/bendeco/GoHub/issues/<number>/comments"
	API_NOTIFICATIONS = API_URL + "/notifications" + ACCESS
)

var activeIssue string
var activeUrl string

var color = map[string] string {
    "Black" : "\033[30m",
    "Red" : "\033[31m",
    "Green" : "\033[32m",
    "Yellow" : "\033[33m",
    "Blue" : "\033[34m",
    "Magenta" : "\033[35m",
    "Cyan" : "\033[36m",

    "Reset" : "\033[0m",
}

func input() string {
	fmt.Println("Enter Action (h, c, q, t)")

	var action string
	fmt.Scanf("%s", &action)

	return action
}

func next(userAction string) {

	switch userAction {
	case "h":
		fmt.Println("Eventually There will be a help message here")
	case "q":
		fmt.Println("YOU CAN NEVER QUIT")
	case "c":
		comment(activeUrl+"/comments"+ACCESS)
		return
	}

	userAction = input()
	next(userAction)
}

func comment(reqUrl string) {
	fmt.Println("Enter Comment Body")

	buf := bufio.NewReader(os.Stdin)
	line, _ := buf.ReadString('\n')

	var commentUrl = strings.Replace(reqUrl, "<number>", activeIssue, -1)

	str := []string{"{\"body\": \"", strings.Replace(line, "\n", "\\n", -1), "\"}"}

	var x = strings.Join(str, "")
	c := bytes.NewBufferString(x)

	res, err := http.Post(commentUrl, "application/json", c)

	if err != nil {
		fmt.Println("http error:", err)
	}

	fmt.Printf("%sPost Successful%s\n\n", color["Green"], color["Reset"])
    fmt.Println(res, reqUrl)
}


// Get all of the comments for an issue
func comments(num, reqUrl string) {
	res, err := http.Get(reqUrl)

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

    fmt.Printf("\n\n")
	for _, val := range c {
        fmt.Printf("-- -- -- -- -- -- -- -- --\n")
		fmt.Printf("%s%s%s -- \n", color["Yellow"], val.User.Login, color["Reset"])
		fmt.Printf("\t%s\n\n", val.Body)
	}

	userAction := input()
	next(userAction)
}

func issue(num, reqUrl string) {
	res, err := http.Get(reqUrl)

	if err != nil {
		fmt.Println("Fatal Error on ISSUE", err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	type Message struct {
		Html_url string
		Url      string
		Title    string
		Body     string
	}

	var m Message
	err = json.Unmarshal(resp, &m)
	if err != nil {
		fmt.Println("Fatal Error", err)
	}

	activeUrl = m.Url

    fmt.Printf("\n%s%s%s (%s)\n\n%s", color["Green"], m.Title, color["Reset"], m.Html_url, m.Body)
	comments(num, API_ISSUE + num + "/comments" + ACCESS)
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

func notif() {
	res, _ := http.Get(API_NOTIFICATIONS)

	resp, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	type Message struct {
		Id      int
		Name    string
		Url     string
		Subject struct {
			Title string
            Url   string
		}
	}

	var m []Message
	_ = json.Unmarshal(resp, &m)

    var notifs = make([]string, len(m))

	for index, val := range m {
        fmt.Printf("%s#%d > %s%s\n", color["Yellow"], index, color["Reset"], val.Subject.Title)
        notifs[index] = val.Subject.Url
	}

	userAction := input()
    if userAction == "c" {
        var num int
        fmt.Scanf("%d", &num)
        thread(num, notifs)
    }
	
}

func thread(num int, notifs []string) {
    /* BODY */
	res, err := http.Get(notifs[num] + ACCESS)

	if err != nil {
		fmt.Println("Fatal Error on ISSUE", err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	type Message struct {
		Html_url string
		Url      string
		Title    string
		Body     string
	}

	var m Message
	err = json.Unmarshal(resp, &m)
	if err != nil {
		fmt.Println("Fatal Error", err)
	}

	activeUrl = m.Url

    fmt.Printf("\n%s%s%s (%s)\n\n%s", color["Green"], m.Title, color["Reset"], m.Html_url, m.Body)


    /* COMMENTS */
    res, _ = http.Get(notifs[num] + "/comments" + ACCESS)

    resp, _ = ioutil.ReadAll(res.Body)
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

    fmt.Printf("\n\n")
	for _, val := range c {
        fmt.Printf("-- -- -- -- -- -- -- -- --\n")
		fmt.Printf("%s%s%s -- \n", color["Yellow"], val.User.Login, color["Reset"])
		fmt.Printf("\t%s\n\n", val.Body)
	}

    fmt.Printf("Action? (c/y/n) > ")
    var answer string
    fmt.Scanf("%s", &answer)

    if answer == "y" {
        // mark read
    }

    if answer == "c" {
        // comment
        comment(notifs[num] + "/comments" + ACCESS)
    }
}

func execute(cmd, num string) {
	switch cmd {
	case "issues":
		issues()
		return
	case "issue":
		issue(num, API_ISSUE + num + ACCESS)
		return
	case "comments":
        comments(num, API_ISSUE + num + "/comments" + ACCESS)
		return
	case "notif":
		notif()
		return
	}
}

func main() {
	var cmd = flag.String("cmd", "", "Help")
	var issueNum = flag.String("num", "", "Help")
	flag.Parse()

	activeIssue = *issueNum
	fmt.Println(activeIssue)

	execute(*cmd, *issueNum)
}
