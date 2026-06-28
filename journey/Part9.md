# Part 9 - Dockerizing and Actions

Now to work on Dockerizing this project for a nice v0.1 release.

Here's how I see it working:

1.  Using a `makefile`, outline specific actions I want to use regularly
2.  When I make a PR, `Github actions` start a job to build an image, and push it to Dockerhub using the `SHA` as the tag
3.  When I merge to `main`, this build is tagged as `latest`

Should be pretty straightforward to start, but in the future I will likely tweak this so that I can create "releases" instead of every build just becoming the latest release. Future tweak though.

Let's start.

First is the `makefile` I want this to contain the following actions:

- `build`
- - build the `go` project
- `build-image`
- - build the docker image
- `build-image-login`
- - handle dockerhub login
- `build-image-push`
- - push the image to dockerhub
- `run`
- - run the `go` project locally

I needed some help with it from AI (boo :( ) but I was able to get a decent makefile setup. Now to work on some actions.

I want the `github action` to do a couple things:

On new PR's:

- Confirm the project builds
- Eventually I will add in running tests, and any FE related building/linting as well

On new releases:

- Build the docker image from the `main` branch
- Tag the image with the release number (ie 1.0, 1.1, etc)
- Tag the image as `latest`
- Push the image to dockerhub
