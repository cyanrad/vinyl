package ingestion

import (
	"context"
	"database/sql"
	"main/db"
	"strings"
)

func (e *Engine) CreateArtists(artists []ArtistIngestion) error {
	for _, artist := range artists {
		description := sql.NullString{String: "", Valid: false}
		if artist.Description != nil {
			description = sql.NullString{String: *artist.Description, Valid: true}
		}

		err := e.queries.CreateArtist(context.Background(), db.CreateArtistParams{
			Name:        artist.Name,
			Description: description,
			Links:       sql.NullString{String: artist.Links.ToString(), Valid: true},
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

		description := sql.NullString{String: "", Valid: false}
		if track.Description != nil {
			description = sql.NullString{String: *track.Description, Valid: true}
		}

		track, err := e.queries.CreateTrack(context.Background(), db.CreateTrackParams{
			Title:       track.Title,
			Description: description,
			Tags:        sql.NullString{String: strings.Join(track.Tags, ", "), Valid: true},
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
