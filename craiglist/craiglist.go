package craigstlist

import (
	"io/ioutil"
	"log"
	"net/http"
 )

 func getPage() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}

	//Convert the body to type string
	sb := string(body)
	log.Printf(sb)
 }

 func main() {
	 getPage()
 }
