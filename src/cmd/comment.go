package comment

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
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
