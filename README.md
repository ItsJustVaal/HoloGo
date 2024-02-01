## HoloGo

# About

The standard youtube interface is not organized, navigatable and seems to never be full up to date
with the live streams on their site. I am putting together a personalized dashboard that utilizes
the youtube API to show specific channels (these channels can be a list of any channels, im using Hololive in this case)

The idea is to make a semi-customizable interface that is easily navigated, well organized and brings together
outside data for each channel to make a one stop shop space that is always up to date.

# Current Technologies Used

Go, Chi, SQLc, PostgreSQL, Goose, HTMX

# TODO

Set up main routes to being building the Front-End
Setup Front-end using HTMX and Go Templates
Add dockerfile & Deploy / Host from my Raspberry Pi
Add an Admin Center for personal commands (I.E. adding and removing channels from the DB)
Add user login and auth

# Notes

The Youtube API is using my personal API Key
this results in incomplete video listings as
streams such as membership only streams will
not populate. I will need to look into this in the future
to properly add user customization.
