package ingestion

import (
	"context"
	"main/db"
	"strings"
)

func (e *Engine) CreateArtists(artists []ArtistIngestion) error {
	for _, artist := range artists {

		err := e.queries.CreateArtist(context.Background(), db.CreateArtistParams{
			Name:        artist.Name,
			Description: artist.Description,
			Links:       artist.Links.ToString(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) CreateTracks(tracks []TrackIngestion) error {
	for _, track := range tracks {
		artistIDs := []int64{}
		for _, artist := range track.Artists {
			artistID, err := e.queries.GetArtistByName(context.Background(), artist)
			if err != nil {
				return err
			}

			artistIDs = append(artistIDs, artistID)
		}

		var tags *string
		if len(track.Tags) > 0 {
			joined := strings.Join(track.Tags, ", ")
			tags = &joined
		}

		track, err := e.queries.CreateTrack(context.Background(), db.CreateTrackParams{
			Title:       track.Title,
			Description: track.Description,
			Tags:        tags,
		})
		if err != nil {
			return err
		}

		for i, artistID := range artistIDs {
			err = e.queries.CreateTrackArtist(context.Background(), db.CreateTrackArtistParams{
				TrackID:  track.ID,
				ArtistID: artistID,
				Rank:     int64(i + 1),
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
