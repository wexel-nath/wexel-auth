# wexel-auth [![Build Status](https://travis-ci.com/wexel-nath/wexel-auth.svg?branch=master)](https://travis-ci.com/wexel-nath/wexel-auth)
Authentication Service using JWT

# Running locally
Set up the database
```
CREATE DATABASE auth;
```
Create tables using the schema .sql files in db/schema.

Start the server from the command line
```
PORT=3000 DATABASE_URL=postgresql://user:secret@localhost:5432/auth go run api/main.go
```
The server should be accessible at localhost:3000/healthz
