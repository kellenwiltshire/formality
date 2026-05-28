# Part 4 - Form routes

Form routes are up next, and they are going to be trickier. To recap, these are the routes I need to build:

`/forms/{id}`

- `GET` Return form information
- `PUT` Update form information
- `DELETE` Delete form
- `POST` Create a new form response - This will be a multi-step process, and may use Queues to handle it all. For now, I will just focus on getting a form response into the DB

`/forms`

- `GET` Return all forms for that user
- `POST` Create a new form for that user

`/forms/{id}/responses`

- `GET` Return responses for specific form

`/forms/{id}/responses/{id}`

- `GET` Return specific form response
- `DELETE` Delete specific form response

These routes will be bare-bones to start. Eventually there will need to be some auth checks so users can only access their forms and their responses, but I can build on that later as part of the total auth part. Building the "Get All" route first to return all forms in the DB was pretty straight forward, it is similar to the get all users route, except now I just need to match on the provided user id. For now, I will pass the users `id` as a query param.

Creating the form will be similar to creating a user as well, just have to pass the items in the body of a POST request. The `user_id` needs to match a valid user to be created, so eventually this will also be handled with auth, but for now we'll pass the `user_id` in the body of the post request. Boom, form created and our previous "Get All" route returns it. Easy, right? Except, because the form `id` is generated to be a 6 character alpha-numeric id, we need to account for the small 1 in 56 Billion possibility that it will generate the same `id` twice, which can't work. This 6 character id is eventually what will be used to determine what responses belong to each form, and thus where to email them. This `id` will be available client side, so it can't be a UUID as that would just be too long.

One way to do this is to just try again, the likelihood of it failing twice in a row for a collision is ridiculously small. This will be an okay solution for now, but if this scales, we'll need a better solution.

Creating the rest of the form routes should be pretty straightforward. Onto the responses...

First off, is receiving a response. For now, we will just load it into the DB. The responses are posted to the `/forms/{id}` route. We then take the form `id` and use this to create an entry in the `form_responses` table, linking them. Eventually, this will also send off an email to that forms `target_email`, but I will work out smtp later.

How this will work though is that the response will be placed into its table with the status of `received`, and quickly the client is notified that the form was submitted successfully. This should be very fast for the client. After, we will kick off the actual email part, that way if it runs into issues, the client won't know. I am thinking of using a Queue for this, but another option may be a cron job that checks for `received` forms, and then processes them. This wouldn't scale well though... A queue would be the right approach, less chance for things to back up, and a better way to handle retry backoffs. That comes later though, for now I have created all the other routes for getting submissions, deleting them, and creating them. Pretty straight forward. Good for this part. Next, we'll create the routes for SMTP settings.
