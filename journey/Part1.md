# Formality

---

### Why am I making this?

I am making this for a couple reasons. First, is I want to move off of using Formspree for my client sites contact forms. It's not convenient, and it can cost money. Secondly, I want to build it to learn more about system design, backend development, databases, and CI/CD deployment with Docker. I want to challenge myself to build something new, even if just for myself. As part of this, I will not use AI to build this outside of a few key situations: 1) To write tests; 2) To have some concepts or documentation explained further. Outside of these reasons, I do not want to use AI at all, especially for code generation.

---

## Functional Requirements

### Frontend

- Users can create forms, but not actual forms, just form-ids.
- - These form-ids can be imported to their application so when a form is submitted this id is used to identify the user and the email the form info should be sent to
- Users can view a history of their form submissions and remove them if they'd like
- Users can update their SMTP settings
- Admin can add other users
- Frontend is only used to manage users, forms, and SMTP settings, will have no interaction with any apps that are using the API services

### Backend

- Spam prevention is on the app side (reCaptcha or other), with ratelimiting on the BE
- Users can delete form submissions
- Users can set what smtp information for mailing form submissions
- Each form can only be sent to 1 recipient
- Authentication handled by session based JWT
- Multiple users allowed on a single service

## Non-Functional Requirements

- Self hosted
- - Should be able to run get setup with a simple docker-compose
- - Frontend UI, Backend, and Database
- Emails should be sent to the recipient within a reasonable timeframe from submission, any lag should be on the smtp provider side
- Once confirmation that an email has been sent, update status in DB
- If email sent failed, try again with fall off strategy (try again right away, then in 30 seconds, than in 1 minute)
- - Might need a queue system for this, with alerts
- - After N fails, the system should mark in the DB that email failed. Cron job to try again for failed emails every X hours
- - If 400 error received from smtp, then stop emails until fixed

## Technology

- NextJS (Frontend UI)
- Go (Backend)
- Queue (RabbitMQ)
- Postgres (Database)

This should provide scale-ability if required, for larger setups. I am aiming for high-availability over consistency. It is better for a User to get an email form sent to them twice than for the customer to experience an error.

Let's get started...
