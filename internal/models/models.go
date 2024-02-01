package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type JsonErrorResponse struct {
	Error string `json:"error"`
}
type JsonResponse struct {
	Status string `json:"status"`
}


// Struct used for each Channel with added
// fields for dashboard personalization options to be added
// in the future
type Channel struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
	Channel   string         `json:"channel"`
	Channelid string         `json:"channelid"`
	Region    string         `json:"region"`
	Prio      sql.NullBool   `json:"prio"`
	Oshi      sql.NullBool   `json:"oshi"`
	Gen       sql.NullInt32  `json:"gen"`
	Tags      sql.NullString `json:"tags"`
	Company   sql.NullString `json:"company"`
}

// Used with the CSV command, will end up being removed when
// the function is moved
type CSVChannel struct {
	Channel   string `csv:"Channel"`
	Channelid string `csv:"ChannelID"`
	Region    string `csv:"Region"`
	Prio      bool   `csv:"Prio"`
	Oshi      bool   `csv:"Oshi"`
	Gen       int32  `csv:"Gen"`
	Tags      string `csv:"Tags"`
	Company   string `csv:"Company"`
}
