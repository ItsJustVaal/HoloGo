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

The cache mostly works, I cannot for the life of me
figure out how youtube is ordering the call returns

To be honest the cache is kinda mostly pointless
until I add an update function too to update the
time fields
