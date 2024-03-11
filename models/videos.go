package models

// Setup the model for videos being passed
// into a template
type VideoTPL struct {
}

// Make funcs to get all videos
// Probably split into different funcs:
// One will get all videos
// One will get region specific videos
// One will get gen specific videos
// One will get individuals videos
// One will get oshi videos
// One will get prio videos
// It will need to get input from the frontend and return the specific videos
// All of them will just return a slice of VideoTPL objects that the tpl will render
