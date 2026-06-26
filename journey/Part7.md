# Part 7 - Re-Organizing and Authentication

It's time to secure my routes, but I want to make some changes. I've decided I want the Go side of things to serve everything, even the FE (will still be react). So I need to change the file structure a bit, to make it a bit more clear. There will no longer be a backend directory, but the `main.go` file will like at the top level. Then I'll have a folder for packages. Eventually I will add in a `Makefile` as well to take care of building and eventual dockerizing.

Then, it is on to Authentication!

For now, I just want to use a simple login form to authenticate. I may add something OIDC in the future, but we'll stick with the basics for now. I'll then use an auth middleware to control the routes and store the session information as a cookie.

Authentication is really as easy as adding 2 middleware functions. One to check if the user has been authenticated, one to check if they are an admin.

```

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			log.Printf("Error getting session: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		iemail, ok := session.Values["email"]
		if !ok || iemail == nil {
			log.Printf("No email in session")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		email, ok := iemail.(string)
		if !ok || email == "" {
			log.Printf("Invalid email in session")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		userId := session.Values["id"].(string)

		values := SessionContext{
			UserID:	userId,
			Email: email,

		}

		// Add the email to the request context
		ctx := context.WithValue(r.Context(), UserKey, values)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func AdminAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			response.HttpResponse(w, "", 0, "Unauthorized", http.StatusUnauthorized)
			return
		}

		email, ok := session.Values["email"].(string)
		if !ok || email == "" {
			response.HttpResponse(w, "", 0, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var count int
		err = database.Db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1 AND role = 'admin'", email).Scan(&count)
		if err != nil || count == 0 {
			response.HttpResponse(w, "", 0, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

```

For the Admin, I want to return the Unauthorized response with our custom `response`, but for the regular auth check, I want to redirect them back to the homepage. This is helpful if they have someone got into a state where they shouldn't be and allows them to log in again.

The `LoginHandler` is where they are authenticated and the session is initiated. This is basic for now, but I'll update with better logging and make sure to set a session length on the cookie later.

```

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login Login
	_ = json.NewDecoder(r.Body).Decode(&login)

	// Get user from DB and check password
	var user User
	var password_hash string
	dbErr := database.Db.QueryRow("SELECT id, role, password_hash FROM users WHERE email = $1", login.Email).Scan(&user.ID, &user.Role, &password_hash)
	if dbErr != nil {
		fmt.Println(dbErr)
		response.HttpResponse(w, "", 0, "Unable to get user", 500)
		return
	}

	// Check the password from the body to the stored hash
	pass_err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(login.Password))
	if pass_err != nil {
		fmt.Println(pass_err)
		response.HttpResponse(w, "", 0, "Unable to login", 500)
		return
	}

	// Get the session
	session, _ := store.Get(r, "session-name")
	// Set session values
	session.Values["email"] = login.Email
	session.Values["role"] = user.Role
	session.Values["id"] = user.ID
	err := session.Save(r, w)
	if err != nil {
		response.HttpResponse(w, "", 0, "Failed to save session: "+err.Error(), 500)
		return
	}

	response.HttpResponse(w, "success", 1, "", 200)
}
```

Overall it is working as expected. I also went ahead and replaced `gorilla/mux` for the routing with `chi` as `gorilla/mux` is no longer maintained. I'm still thinking there is a lot left to be desired for how my project is laid out. Going take some more time to get it to a place I feel a little better about it.

## Many, Many days later

I should have committed what I had before I "took some time" to get things a little more laid out. In that time, I completed a course by `Melkeydev` on Frontend Masters around Go and it really opened up a lot for me. I completely re-organized the entire project structure based on this new knowledge. The new folder structure

```
- /internal // This is for everything that could be considered "backend"
| - /app // This handles the main "brain" of the Go project. It contains all the connections between different parts of the service
| - /api // This contains all the "handler" files for the project. A "handler" is anything that interacts with a route
| - /listener // This contains the DB listener in order to set off the email part of the service
| - /middleware // This contains the middleware, obviously. For now its main purpose to handling authentication checks
| - /routes // This contains all the routing logic for the application. This connects to the handlers
| - /service // This contains any service related code. For now, it only contains the logic around sending emails. It connects to the listener and one route for testing
| - /store // This contains all the "store" files for the project. A "store" is anything that interacts with the database
| - /tokens // This contains the token issuing logic. For now, we hand the tokens back to the client for them to be part of the request as a "Bearer" token, but I will be updating this in the next part to use cookies instead
| - /util // This contains any utility files that are needed across multiple different parts of the application
- /migrations // Contains all the DB migration files
```

Basically how it all ties together is:

```
User request -> route -> handler -> store -> service/token // Writing this out I think I should maybe move tokens into service, we'll see.
```

It was a big overhaul, but I think the project is in a lot cleaner, and clearer, of a state to continue forward.

Next step is further testing, updating auth to use cookies, and attempting a docker build process I think.
