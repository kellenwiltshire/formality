# Part 8 - Cookies and Testing

So for this part I want to work on updating the authentication to use cookies instead of looking for a `Bearer` token on the request. Since this is running all from the same location, this seems like the right thing to do. If I ever want to allow for cross-origin type of work, I may need to update, but for now I will assume any authenticated requests will be coming from the same domain.

I can re-use a log of the logic I have around creating the tokens, so I don't need to make any drastic changes to the store. Just instead of returning the token to the user, I'll set it in a cookie. I can then check this cookie for the token and confirm it is still alive and well in the DB when a request comes in.

I set this cookie like so:

```
http.SetCookie(w, &http.Cookie{
    Name:     "formality_auth",
    Value:    string(token.Hash),
    Expires:  token.Expiry,
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteLaxMode,
    Path:     "/",
})
```

So it mirrors what I've set in the DB, even the expiry time (currently set to expire 24 hours after issue). So when the cookie is gone, it should also be expired in the DB. I'll need to update the logout handling, currently it logs the user out from EVERYWHERE. It should log them out only from their current device. Future problem though.

Now I need to update the middleware to look for the cookie, and not the `Bearer`

```
c, err := r.Cookie("formality_auth")
if err != nil {
    if err == http.ErrNoCookie {
        util.WriteJSON(w, http.StatusUnauthorized, util.Envelope{"error": "token expired or invalid"})
        return
    }
    util.WriteJSON(w, http.StatusInternalServerError, util.Envelope{"error": "internal service error"})
    return
}

token := c.Value
```

Boom, now we set the token to be the value of the cookie, and continue to check the DB and set the user from there. In theory, this works. Now, I've gotta do some serious testing of all my changes. I think I should squash my migrations and start fresh though. It only removes 2 migration files, but might as well keep it clean while I work in a blank environment.
x
Testing the cookie part seems to be tough without a browser to store the cookies... Luckily `Bruno` has a client that works on Linux!

Now I just work my way through the routes and test they work as expected... Boring, but required.

With a bit of work, and some slight tweaks, all routes work as expected and smtp email works. Success! I thinK I can dockerize this part and it can function as a v0.1. There is no UI, and the email it sends is pretty ugly, but it works as a form handler. Not ready for release, but ready to start incrementing on towards a v1.

Dockerizing should be pretty straight-forward, I hope. I also want to create a `makefile`, and some `github actions` to ensure this all happens pretty seamlessly.

Actually, this will be in the next part...
