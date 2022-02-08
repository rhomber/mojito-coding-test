# SOLUTION

This document attempts to explain the choices made in providing this solution.

# common

The first thing you may have noticed is a package called __common__, this is where I have placed all the common
libraries. I usually build micro-service based architectures which pull in common resources (usually included using
a __private__ directory with git submodules). So the files within __common__ are essentially a cut down version of 
the __core__ library I use in my personal projects.

## boot

A means of booting a go project, the purpose is to reduce boiler-plate code in the projects _main.go_ file. I also
find this really useful to ensure certain things are always done (i.e. allocating enough resources).

## chttp

Extended (chi) functionality for HTTP request handling. I primarily do this as a means to reduce repeated boiler-plate
code for passing resources around.

## config

Extended (viper) functionality for reading a config file as well as environment variables. I like to prefix my
environment variables with GO_APP_ to avoid collisions (in the same way that React App does).

## data

High level data structures / types.

## errs

A simple error framework which is used to expose special error codes and rich data structures. When I first implemented
this Go didn't have error wrapping (or perhaps I wasn't aware of it). I would probably have designed this differently
if I were to start with a clean slate. However, I feel exposing error codes is very important if you want the UI
experience to be better for the end user. I also like to expose _correlation id's_ in the user displayed error messages
as well. The idea being a _correlation id_ is carried throughout a request, through all micro-services. So when you scan 
the logs you can see everything that happened.

