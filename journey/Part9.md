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
- Push the image(s) to dockerhub

This way there is only a new version for users to update to when there is a new release. I don't need a CI/CD setup for this, unlike a website, since the updates wouldn't be that often since it is self-hosted.

I've merged my changes so that I can have the actions ready. Now to test them. I'll create a new PR with this update journey log, and if that starts an image build, I'll know one of them is correct! However, I am thinking that the build is likely to fail, since I am not sure if a build will succeed without environment variables. It should, but we'll see.

It passed! Merged. Not to try a release. A brand new release of v0.0.1!

It failed. Issue with the login process for dockerhub. Hmm...

I think I may be passing flags wrong to the makefile commands, if this doesn't work, I may try another way to login to Dockerhub in actions. Of course, I think I ALSO need to create a new release to test any changes...

I got it to work! Damn typos... But on release it build an image, then tagged the image with both the release tag, and the latest tag, and then uploaded them to docker. Perfect.

Good for this part. Next up, integrating with `docker-compose` which means passing `env` variables into the container, rather than creating an `env` file. Will be a big change, but hopefully not too bad.

Which was also pretty easy! Use a local build context for the image (going to try pulling the image too once I have my tweaks in) and then update to pull environment variables from the environment and not a file. Pretty straight-forward. Up and running like new. I also setup a `Dockerfile.dev` that I can run with a new tool I found called `Air` which watches for my go change. No need to keep re-loading every time I make a change. Awesome.
