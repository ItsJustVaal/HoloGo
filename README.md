# HoloGo

A back end API that uses Youtube API to track and store videos
for all holomems. Stores them in a postgres db for frontend use.

# Technologies

Go, Chi, sqlc, Postgresql, Goose, HTMX

# TODO

Add dockerfile to host on my pi later
Set up API and routes  
Add basic frontend (going to try HTMX)

# Notes

Fixed the cache by adding the published at date converted
to RFC3339 and ordered by that, and fixed the ticker loop that was causing
the goroutine to call the same playlist ids sometimes
for some reason.

Added a github workflow for formatting checks

Next is to set up a few basic routes
then a basic front end
