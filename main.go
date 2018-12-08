package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// User ... holds the data of users in system
type User struct {
	ID               int      `json:"id"`
	UserName         string   `json:"user_Name"`
	Ips              []string `json:"ips"`
	Target           string   `json:"target"`
	EVENT0ACTION     string   `json:"EVENT_0_ACTION"`
	DateTimeAndStuff int      `json:"DateTimeAndStuff"`
}

func main() {

	var token string    // Authentication token
	var users []User    // Will hold the users data that is read from the API
	var fileName string // Name of JSON file to store data

	// Get the token
	token = getToken()

	// Make a GET request to /get-events to get all users data
	users = getUsersAPI(token)
	fileName = "usersData.json"
	writeToFile(users, fileName)

	/* When the program runs the web app can be viewed in the browser by going to
	127.0.0.1:8000/users
	*/
	http.HandleFunc("/", usersPageHandler)
	http.HandleFunc("/users", usersPageHandler)
	http.ListenAndServe(":8000", nil)

} // end func main

//***************************************************************************************************

/* ... This function makes a GET request to /auth and returns a Authentication
token.
*/
func getToken() string {
	//make a client
	client := &http.Client{}

	res, err := client.Get("https://duoauth.me/auth")
	if err != nil {
		fmt.Printf("The HTTP failed with errors %s\n", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("An error when trying to read the response body %s\n", err)
	}
	defer res.Body.Close()
	return string(body)
}

/* ... This function makes a GET request to /get-events and returns users data
 */
func getUsersAPI(tok string) []User {

	client := &http.Client{}
	var userData []User

	resp, err := http.NewRequest("GET", "https://duoauth.me/get-events?from=1&to=2000", nil)
	resp.Header.Set("Authorization", tok)
	resps, err := client.Do(resp)
	if err != nil {
		fmt.Printf("The HTTP failed with errors %s\n", err)
		fmt.Println(err.Error())
	}

	data, err := ioutil.ReadAll(resps.Body)
	if err != nil {
		fmt.Printf("The HTTP failed with errors %s\n", err)
		fmt.Println(err.Error())
	}

	// Decode JSON data into User struct
	err2 := json.Unmarshal(data, &userData)
	if err2 != nil {
		fmt.Printf("Error JSON Unmarshalling %s\n", err)
		fmt.Println(err.Error())
	}
	defer resps.Body.Close()

	return userData
}

/*... This function encodes the users data and copies it to a JSON file.
 */
func writeToFile(u []User, fileName string) {

	// Endcode the data and copy it to  users.jason
	var buf = new(bytes.Buffer)

	enc := json.NewEncoder(buf)
	enc.Encode(u)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error JSON marshalling %s\n", err)
	}
	io.Copy(file, buf)
	defer file.Close()

}

/*... The function handles the /users request that dispalys the users data
 */
func usersPageHandler(w http.ResponseWriter, r *http.Request) {
	var templateUsers []User

	templateUsers = readFromFile("usersData.json")
	t, err := template.ParseFiles("users.html")
	if err != nil {
		fmt.Printf("The file users.html was not found %s\n", err)
	}
	t.Execute(w, templateUsers)
}

/*... This function decodes the users data from a JSON file
 */
func readFromFile(fileName string) []User {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("The file users.json was not found %s\n", err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	var usersDecoded []User
	dec.Decode(&usersDecoded)

	return usersDecoded
}
