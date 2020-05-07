# go-web 
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cbigge/go-web)

This web project is a image gallery hosting platform that I built to learn web programming in golang. It is built using a MVC architecture. The `models/` folder contains all the data models and database communication. The `controllers/` folder contains all communication logic between the models and views, as well as manages routes. The `view/` folder contains all frontend pages and templates.
There are other utility folders in the project that serve different functions. 
- The `context/` folder contains functionality to store the current user in the web server's context, making it accessable from the controller. 
- The `hash/` folder contains functionality to secure user credentials in the database through a hmac hash encryption. 
- The `middleware/` folder contains functionality that checks user authentication for access to view specific routes. 
- The `rand/` folder contains functionality for string encoding that is used for a user's 'remember_token' which is how the server knows if they are authenticated.

## Prerequisites:
- Go 1.13+
- Postgres

## Before you run:
- Create a database in postgres
- Update the config.json file with the database information and postgres credentials

## How to run:
- Compile and start by running `go run *.go` in the base project directory. The `*.go` is required because there is two files in the base project directory that need to be compiled and ran together: `main.go` and `config.go`

## Usage:
- First, you need to create an account; you can't do anything without one.
- After signing up, you will see that a new item is on the navigation bar titled `Galleries`. Go ahead and navigate to it.
- On the galleries page, you will be able to see all of the current galleries, as well as create your own.