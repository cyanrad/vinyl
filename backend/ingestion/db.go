package ingestion

import (
	"database/sql"
	"log"
	"main/db"
	"main/ingestion/storage"
	"main/util"
	"strings"
)

// TODO: in case an artist is mentioned that can't be traced back it show log the error but continue
func (e *Engine) CreateTracks(tracks []storage.TrackIngestion) error {
	for _, track := range tracks {
		// checking if track already exists
		albumName := track.Album
		if albumName != nil {
			albumName = &strings.Split(*track.Album, " - ")[1] // extracing the album name part
		}
		trackName := util.GenerateTrackName(util.GenerateConcatNames(track.Artists), albumName, track.Title)
		log.Printf("Checking if track %s exists", trackName)
		_, err := e.queries.GetTrackByName(e.ctx, trackName)
		if err == nil {
			log.Println("Track already exists. Skipping db insert")
			continue
		} else if err == sql.ErrNoRows {
			log.Println("New track found, inserting into db")
		} else {
			return err
		}

		// getting artist idea to link to track
		artistIDs := []int64{}
		for _, artist := range track.Artists {
			artistID, err := e.queries.GetArtistByName(e.ctx, artist)
			if err != nil {
				return err
			}

			artistIDs = append(artistIDs, artistID)
		}

		// getting album id to link
		var albumID int64 = -1
		if track.Album != nil && track.AlbumRank != nil {
			id, err := e.queries.GetAlbumByName(e.ctx, *track.Album)
			if err != nil {
				return err
			}

			albumID = id
		}

		// generating tags string
		var tags *string
		if len(track.Tags) > 0 {
			joined := strings.Join(track.Tags, ", ")
			tags = &joined
		}

		// creating the track
		trackRow, err := e.queries.CreateTrack(e.ctx, db.CreateTrackParams{
			Title:       track.Title,
			Description: track.Description,
			Tags:        tags,
		})
		if err != nil {
			return err
		}

		// linking the artists
		for i, artistID := range artistIDs {
			err = e.queries.CreateTrackArtist(e.ctx, db.CreateTrackArtistParams{
				TrackID:  trackRow.ID,
				ArtistID: artistID,
				Rank:     int64(i + 1),
			})
			if err != nil {
				return err
			}
		}

		// linking the album
		if albumID != -1 {
			err = e.queries.CreateTrackAlbum(e.ctx, db.CreateTrackAlbumParams{
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

func (e *Engine) CreateAlbums(albums []storage.AlbumIngestion) error {
	for _, album := range albums {
		albumName := util.GenerateAlbumName(util.GenerateConcatNames(album.Artists), album.Name)
		log.Printf("Checking if album %s exists", albumName)
		_, err := e.queries.GetAlbumByName(e.ctx, albumName)
		if err == nil {
			log.Println("Album already exists. Skipping db insert")
			continue
		} else if err == sql.ErrNoRows {
			log.Println("New album found, inserting into db")
		} else {
			return err
		}

		artistIDs := []int64{}
		for _, artist := range album.Artists {
			artistID, err := e.queries.GetArtistByName(e.ctx, artist)
			if err != nil {
				return err
			}

			artistIDs = append(artistIDs, artistID)
		}

		albumRow, err := e.queries.CreateAlbum(e.ctx, db.CreateAlbumParams{
			Name:        album.Name,
			Description: album.Description,
		})
		if err != nil {
			return err
		}

		for _, artistID := range artistIDs {
			err = e.queries.CreateArtistAlbum(e.ctx, db.CreateArtistAlbumParams{
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

func (e *Engine) CreateArtists(artists []storage.ArtistIngestion) error {
	for _, artist := range artists {
		log.Printf("Checking if artist %s exists", artist.Name)
		_, err := e.queries.GetArtistByName(e.ctx, artist.Name)
		if err == nil {
			log.Println("Artists already exists. Skipping db insert")
			continue
		} else if err == sql.ErrNoRows {
			log.Println("New artist found, inserting into db")
		} else {
			return err
		}

		err = e.queries.CreateArtist(e.ctx, db.CreateArtistParams{
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
