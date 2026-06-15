# Part 7 - Re-Organizing and Authentication

It's time to secure my routes, but I want to make some changes. I've decided I want the Go side of things to serve everything, even the FE (will still be react). So I need to change the file structure a bit, to make it a bit more clear. There will no longer be a backend directory, but the `main.go` file will like at the top level. Then I'll have a folder for packages. Eventually I will add in a `Makefile` as well to take care of building and eventual dockerizing.

Then, it is on to Authentication!
