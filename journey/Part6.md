# Part 6 - Form Submission pipeline

Now that all the routes are defined, and SMTP is working. It's time to start simulating and working through someone submitting a form response.

This process is fairly straightfoward. A Customer submits a form on our Users website, and the website makes a `POST` request to `/forms/{id}` where `id` is the id of the form that the User has created. Once that API receives the request, it stores the form response in the DB and will then email the contents of the form to the User, updating the DB based on its success. For now I will not worry about retrying on a failed send, or any type of notification to the user on a failed send, that'll come later.

For testing, the `body` of the `POST` request will look like this:

```
{
    "name": "Test Name",
    "email": "test@email.com",
    "message": "Test form submission message from a customer"
}
```

I've already got the part built that saves this data into the DB, so I really just need to add in the next part of that function that calls the SMTP pipeline. I'll just need to grab the `user_id` associated with that form. I had to make some changes about the response is saved in the DB, to return the `response_id`. Everything looked good, until I got the email and the body simply contained

```
Empty Message
```

Must be something wrong with getting the payload and putting it into the message. Time to investigate.

The issue was a typo, of course... `"SELECT payload FROM form_submissions WHERE id = $2"`... `$2` should be `$1` of course.

The email body now reads:

`{"name": "Test Name", "email": "test@email.com", "message": "Test form submission message from a customer"}`

Not pretty, but a start! I'll have to build a html body builder for this later. The issue though is it took over 5 seconds for the Customer to get a 200 response after submitting. That needs to be faster. It should be nearly instant! I want it to return the 200 to the user as soon as we insert the submission into the DB, not wait until the email has been sent. I think I need to change my approach and implement a listener that can see when a Table gets a new entry, I think something like `pg_notify` is what I want, or at least the first thing to try.

First we create the `TRIGGER` and `FUNCTION` in the DB:

```

CREATE OR REPLACE FUNCTION notify_form_submissions_insert()
RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify(
        'form_submissions_inserts',
        json_build_object(
            form_id: NEW.form_id,
            payload: NEW.payload
        )::text
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER form_submissions_insert_trigger
AFTER INSERT ON form_submissions
FOR EACH ROW
EXECUTE FUNCTION notify_form_submissions_insert();

```

Then we can create a new package to handle this, then load it into `main()`. Unfortunately, I did his a snag here setting up the listener how I wanted and had to get some help from AI :(. But once the listener is working I just had to lookup (on my own!) how to get it so the DB listener and routes listener can run at the same time. Turns out I can use something called `sync` to run both. Updating the smtp `PrepareMail` function to accept the new parameters was easy, and boom. It works! Now the Client request is super fast (for me, all local, it was around 10ms) and then a few moments later I received the email. Amazing. A bit of tidying up and I think this part is done. Getting close to the end of the backend I think, but still need to get into authentication and security stuff! Fun!
