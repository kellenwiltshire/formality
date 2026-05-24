# Part 3 - The API-ining

Time to start the API building. But first, I need to figure out how to organize in Go. I can't have a 10k line `main.go`. I am going to move the migration stuff into the migration folder, into a `migration.go` file and import it into `main.go`. Took a bit to learn the syntax but now, my `main.go` just imports the migration function and we're good to go. Now to build out the api. I'll do this for now in the `main.go` file, then move it out into its own file structure later. Let's get Bruno up and running and try to get some some data returned!

```
func hello(w http.ResponseWriter, r *http.Request) {
	 w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode("Hello!")
}


func main() {
	migration.Migration()

	r := mux.NewRouter()

	r.HandleFunc("/hello", hello).Methods("GET")

	//Start the Server
	port := goDotEnvVariable("PORT")
	fmt.Println("The server is running on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
```

It works! "Hello!" is returned on that endpoint. Now we can get cooking here. Let's map out the initial endpoints and what they'll do.

#### User Endpoints

`/users/{id}`

- `GET` Return user information
- `PUT` Update user information
- `DELETE` Delete user

`/users`

- `GET` Return all users information
- `POST` Create a new user

#### Form Endpoints

`/forms/{id}`

- `GET` Return form information
- `PUT` Update form information
- `DELETE` Delete form
- `POST` Create a new form response

`/forms`

- `GET` Return all forms
- `POST` Create a new form

`/forms/{id}/responses`

- `GET` Return responses for specific form

`/forms/{id}/responses/{id}`

- `GET` Return specific form response
- `DELETE` Delete specific form response

#### SMTP

`/email-settings`

- `GET` Return SMTP settings
- `POST` Create SMTP settings
- `PUT` Update SMTP settings
- `DELETE` Delete SMTP settings

This seems okay to start, will probably see a lot of changes once I get building. I'll start with creating a `routes` package in go to handle all these routes. I will also need to create route protection, but that can come later. First, I'll start with the user routes and see if I can add, get, update, and delete users in my database.

---

That was a learning experience. I think the way I have done this is very inefficient, and will need a lot of updating, but I am now able to add, update, get, and delete users. So that's a win. I think I am going to have to make some changes around my `structs`, and I want to implement some sort of standardized responses. My thinking is, no matter what, every response from the API should have the following shape:

```
{
    data: // This could be either an array of objects, or an object
    status: // This will either be 1 for success, or 0 for error
    message: // This will contain any messages I want to pass along. For instance, a successful action without needing to return data can state that,
                and an error can return an error response to pass along to the user (if necessary)
    httpCode: // Pass along an appropriate http status code, default is 200
}
```

For this to work, I think I need to build a standard `http response` function to handle all of this. This should clean up some of the code as well going forward. This is what I am thinking:

```
func HttpResponse [T any](w http.ResponseWriter, data T, success int, message string, status int) {
	var response Response

	response.Data = data
	response.Success = success
	response.Message = message

	if success < 1{
		fmt.Println(message)
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
```

Now when I build out my frontend, I will know exactly what shape every reply will be. User routes are not complete, for now. Auth will come later around those routes. Good spot to leave for now and work on the form routes next.
