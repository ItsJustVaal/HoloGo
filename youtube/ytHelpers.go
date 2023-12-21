package youtube

import (
	"context"
	"fmt"

	"github.com/ItsJustVaal/HoloGo/internal/database"
)

// func startAPICalls(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
// 	log.Printf("Collecting videos every %s on %v goroutines...", timeBetweenRequest, concurrency)
// 	ticker := time.NewTicker(timeBetweenRequest)

// 	for ; ; <-ticker.C {
// 		channels, err := db.GetChannelIDs(context.Background())
// 		if err != nil {
// 			log.Println("Couldn't channel IDs", err)
// 			continue
// 		}
// 		log.Printf("Found %v feeds to fetch!", len(channels))

// 		wg := &sync.WaitGroup{}
// 		for _, ch := range channels {
// 			wg.Add(1)
// 			// go ch start getting ids
// 		}
// 		wg.Wait()
// 	}
// }

func GetChannelIDs(db database.Queries) ([]string, error) {
	channelIds, err := db.GetChannelIDs(context.Background())
	if err != nil {
		return channelIds, fmt.Errorf("failed to get Channels")
	}

	return channelIds, nil
}
