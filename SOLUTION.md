# SOLUTION

This document attempts to explain the choices made in providing this solution.

# General

## README

The __README__ file contains the relevant notes needed to build and run the project.

## Postman Collection

I have saved an export of a __Postman__ collection you can use to test the project under the __testing__ directory.
You will also want to define two environment variables __auctionLotId__ and __authUserEmail__ (an email for a user
you created).

# Structure

## common

The first thing you may have noticed is a package called __common__, this is where I have placed all the common
libraries. I usually build micro-service based architectures which pull in common resources (usually included using
a __private__ directory with git submodules). So the files within __common__ are essentially a cut down version of 
the __core__ library I use in my personal projects.

### boot

A means of booting a go project, the purpose is to reduce boiler-plate code in the projects _main.go_ file. I also
find this really useful to ensure certain things are always done (i.e. allocating enough resources).

### chttp

Extended (chi) functionality for HTTP request handling. I primarily do this as a means to reduce repeated boiler-plate
code for passing resources around.

### config

Extended (viper) functionality for reading a config file as well as environment variables. I like to prefix my
environment variables with __GO_APP___ to avoid collisions (in the same way that React App does).

### data

High level data structures / types.

### errs

A simple error framework which is used to expose special error codes and rich data structures. When I first implemented
this Go didn't have error wrapping (or perhaps I wasn't aware of it). I would probably have designed this differently
if I were to start with a clean slate. However, I feel exposing error codes is very important if you want the UI
experience to be better for the end user. I also like to expose _correlation id's_ in the user displayed error messages
as well. The idea being a _correlation id_ is carried throughout a request, through all micro-services. So when you scan 
the logs you can see everything that happened.

## app

All the files specific to running this application are kept here. __main.go__ loads the application bootstrapping 
code found in __app.go__. I would also prefer you focus on the files here rather than in __common__. These are the 
files I have been working on in order to provide the solution for you. I have included the __common__ files in order
to save a little time, but also demonstrate how I usually structure projects.

### data

The application specific data entities (including models) and types are stored here. I like to use the same approach
as Java with my models, I use DTO's. There are several reasons I prefer to use DTO's, but the primary reason
is for security. You won't accidentally serialise a model and send sensitive information if you're using DTO's. 
Another big reason is for performance, it's often desirable to combine information from several modals together.
To bridge models and DTO's I use adapters, which is also a fairly normal pattern from Java.

### handler

The API handlers are kept here (exposing the view). Here we also store the health check's and any other application
specific middlewares.

### service

The services (or controllers) are kept here. They are invoked from the handlers or other services.

## db

Where the db migrations are stored. I prefer to use a library like __goose__ to use raw SQL migrations than
some ORM auto-creating magic. This is especially important when you want to define exactly how things should be
created (foreign key constraints, unique constraints, potentially triggers). I also feel it's a lot safer, I have
had some bad experiences with the code based generators.

Also, I would like to mention that I usually use PostgreSQL (and CockroachDb), so I prefer to use go-pg (and my 
rewrite of go-pg that supports cockroachdb). This is my first time using sqlite3, and I haven't really used GORM
before either. I would have also used int64 (or uint64) for the ID's as well, but GORM defaults to uint so I went
with it (but would have created my own solution if this were a real scenario).

# Technical Notes

## Users

The requirements didn't mention a user service, but I felt it would be better to provide something simple (for 
better db schema design). See the __Postman__ collection for examples on how to create and query users.

## Authentication

The requirements didn't mention an authentication system, but I took the liberty of implementing a very rustic one.
I felt it better demonstrated how I would lay the project out. This simple authentication system (which doesn't
authenticate) works by setting a header called __X-User-Auth__ to the user's email address. Authentication is only 
enabled on the __/v1/auction/lot/{id}/bid__ endpoints. Usually I would prefer to use something like JWT.

## Auction Lot Bid

### Db Tables

I opted for storing the bids in two different tables vs. one. My reason for doing this was I felt to honor the
requirement of auditing the bids, I wanted the data to be immutable. I also wanted to add a new bid record to the
bid table when the __max bid__ value caused the original bid to be incremented. To make sense of these dynamic bids
I have also introduced a __type__ column to distinguish between __user__ generated bids and __max bid__ generated bids.

## Testing

As the project took considerably longer than I was expecting I have opted for adding unit tests only to one service.
I feel that regardless of time constraints code such as this, with many possible outcomes, should be thoroughly tested.

## Repository

In the past I have also included another layer beyond the service for access to the database (repository). This isn't
what I have been doing lately in my projects as I found the repository and service layers had too much overlap. But
I can see the benefits of having it, if that's what you are using.