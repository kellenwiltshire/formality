# Part 5 - SMTP

This part we will look at saving a users SMTP settings to the Database. Each user can set their own SMTP settings so that they can send the form submissions to their desired location (likely their own inbox).

This is a little trickier though, since the SMTP password should be encrypted. Once it is sent to the BE from the FE, the password will never go back to the FE again. However, when we store it in the DB we should encrypt it and only decrypt it when we need to use it form sending emails.

Since we will need to decrypt the password to use it, we can't use something like bcrypt to encrypt the password, since this is a 1-way encryption method.

We can use `crypto/aes` though. All the Administrator will need to do is generate a 32 byte key for their `.env` and we can use this to encrypt the key. Then we just decrypt it on the server when we need it.

Let's build the routes and try it out.

### The routes

`/email-settings/{user_id}`

- `GET` Return SMTP settings
- `POST` Create SMTP settings
- `PUT` Update SMTP settings
- `DELETE` Delete SMTP settings

When I encrypted the password, it initially has it as a type of `[]byte` which I don't want to save in my DB. It's simple to update it to be a string, though

```
encrypted_password, err := encrypt_text.EncryptAES(smtp_settings.Password)
	if err != nil {
		response.HttpResponse(w, "", 0, err.Error(), 500)
		return
	}

converted_pass := base64.StdEncoding.EncodeToString(encrypted_password)
```

When I need to get the string for sending emails, I just need to do the reverse!

Creating the routes was straightforward, and I am able to encrypt the password before storing it. This is all the routes not completed that I think I need to start, now I can work on the actual sending emails part of it all!

Finding a package to help with sending emails in Go will be the first part. I don't think I need to re-invent the wheel for this project to just send the emails. `go-mail` seems to be pretty highly recommended, but Go has their own package as well `net/smtp`. The built in Go one might be the best for this, I don't need a lot of frills. Let's start with that.

When I call for an email to be sent, I need to pass to that function (simply called `SendMail`) the `user_id` to get the smtp settings, and the response `id` to get the form response from the DB.

An oversight though, I never added a column in the `smtp` table for a recipient address for these emails to go to! Nor did I add a `from` column, I'll have to alter the DB for these.

I'll also need an easy way to test this, without creating a new form every time. I think I need a test route, which can fire off a test email to make sure the smtp settings are hooked up right. This was a good thing to work through. I was missing some DB data, and I needed to re-work how my encryption process went. But, in the end I was able to receive an email from the application! We're really cooking now.

I think this is a good place to end this part. All the routes are defined and SMTP is working properly. Next part will be working through the form submission part!
