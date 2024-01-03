# HoloGo

A back end API that uses Youtube API to track and store videos
for all holomems. Stores them in a postgres db for frontend use.

# Technologies

Go, Chi, sqlc, Postgresql, Goose, HTMX

# TODO

Add docker and github stuff  
Set up API and routes  
Add basic frontend (going to try HTMX)

# Notes

Fixed the cache by adding the published at date converted
to RFC3339 then using that to set the cache.

Next is to set up a basic front end
