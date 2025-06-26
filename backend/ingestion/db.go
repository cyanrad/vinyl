package ingestion

import (
	"context"
	"main/db"
	"strings"
)

// TODO: in case an artist is mentioned that can't be traced back it show log the error but continue

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

		var albumID int64 = -1
		if track.Album != nil && track.AlbumRank != nil {
			id, err := e.queries.GetAlbumByName(context.Background(), *track.Album)
			if err != nil {
				return err
			}

			albumID = id
		}

		var tags *string
		if len(track.Tags) > 0 {
			joined := strings.Join(track.Tags, ", ")
			tags = &joined
		}

		trackRow, err := e.queries.CreateTrack(context.Background(), db.CreateTrackParams{
			Title:       track.Title,
			Description: track.Description,
			Tags:        tags,
		})
		if err != nil {
			return err
		}

		for i, artistID := range artistIDs {
			err = e.queries.CreateTrackArtist(context.Background(), db.CreateTrackArtistParams{
				TrackID:  trackRow.ID,
				ArtistID: artistID,
				Rank:     int64(i + 1),
			})
			if err != nil {
				return err
			}
		}

		if albumID != -1 {
			err = e.queries.CreateTrackAlbum(context.Background(), db.CreateTrackAlbumParams{
				TrackID: trackRow.ID,
				AlbumID: albumID,
				Rank:    int64(*track.AlbumRank),
			})
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (e *Engine) CreateAlbums(albums []AlbumIngestion) error {
	for _, album := range albums {
		artistIDs := []int64{}
		for _, artist := range album.Artists {
			artistID, err := e.queries.GetArtistByName(context.Background(), artist)
			if err != nil {
				return err
			}

			artistIDs = append(artistIDs, artistID)
		}

		albumRow, err := e.queries.CreateAlbum(context.Background(), db.CreateAlbumParams{
			Name:        album.Name,
			Description: album.Description,
		})
		if err != nil {
			return err
		}

		for _, artistID := range artistIDs {
			err = e.queries.CreateArtistAlbum(context.Background(), db.CreateArtistAlbumParams{
				AlbumID:  albumRow.ID,
				ArtistID: artistID,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

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
