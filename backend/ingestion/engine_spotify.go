package ingestion

import (
	"log"
	"main/ingestion/spotify"
	"main/ingestion/storage"
	"main/util"
)

func (e *Engine) IngestSpotify(resourceType util.ResourceType, resourceID string) error {
	util.Log.Debug("Connecting to spotify...")
	spotifyConn, err := spotify.Connect(e.ctx, e.cache)
	if err != nil {
		panic(err)
	}
	util.Log.Debug("Connected to spotify")

	if util.BUST_CACHE {
		util.Log.Infof("Bust cache enabled. Removing cached object for Spotify resource: %s - ID: %s\n", string(resourceType), resourceID)
		err := e.cache.Delete(resourceType, util.SOURCE_SPOTIFY, resourceID)
		if err != nil {
			return err
		}
	}

	util.Log.Debugf("Beginning ingestion of spotify %s %s", resourceType, resourceID)
	switch resourceType {
	case util.PLAYLISTS:
		err = e.IngestSpotifyPlaylist(spotifyConn, resourceID)
	case util.TRACKS:
		err = e.IngestSpotifyTrack(spotifyConn, resourceID)
	default:
		log.Fatalf("ingestion type %s for spotify is not implemented. terminating program", resourceType)
	}

	return err
}

func (e *Engine) IngestSpotifyPlaylist(s *spotify.SpotifyConn, playlistID string) error {
	data, err := s.GetPlaylistData(playlistID)
	if err != nil {
		return err
	}

	ingestion, err := s.GenerateTrackSpotifyIngestion(data)
	if err != nil {
		return err
	}

	err = e.IngestSpotifyData(ingestion)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) IngestSpotifyTrack(s *spotify.SpotifyConn, trackID string) error {
	util.Log.Debug("Fetching track from spotify API")
	data, err := s.GetFullTracks(spotify.StrToID(trackID))
	if err != nil {
		return err
	}

	util.Log.Debug("Transforming track data from spotify to an ingestable form")
	ingestion, err := s.GenerateTrackSpotifyIngestion(data)
	if err != nil {
		return err
	}

	util.Log.Debug("Ingesting the track data and storing it in our database")
	err = e.IngestSpotifyData(ingestion)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) IngestSpotifyData(ingestion spotify.SpotifyIngestion) error {
	// ingesting artists
	for _, artist := range ingestion.Artists {
		// storing data in db
		if err := e.CreateArtists([]storage.ArtistIngestion{
			{Name: artist.Name, Links: artist.Links, Description: nil},
		}); err != nil {
			return err
		}

		// fetching artist image
		if artist.ImageURL != "" {
			storage.FetchImage(util.ARTISTS, artist.ImageURL, artist.Name)
		}
	}

	// ingesting albums
	for _, album := range ingestion.Albums {
		// storing data in db
		if err := e.CreateAlbums([]storage.AlbumIngestion{
			{Name: album.Name, FullName: album.FullName, Artists: album.Artists, Description: nil},
		}); err != nil {
			return err
		}

		// fetching album image
		if album.ImageURL != "" {
			storage.FetchImage(util.ALBUMS, album.ImageURL, album.FullName)
		}
	}

	// ingesting tracks
	for _, track := range ingestion.Tracks {
		// creating track
		var albumRank *int = nil
		if track.Album != nil {
			albumRank = &track.AlbumRank
		}

		if err := e.CreateTracks([]storage.TrackIngestion{
			{
				Title:       track.Title,
				FullTitle:   track.FullTitle,
				Artists:     track.Artists,
				Album:       track.Album,
				AlbumRank:   albumRank,
				Description: nil,
				Tags:        nil,
			},
		}); err != nil {
			return err
		}

		// we fetch the image if the album doesn't exist or is a single
		if track.Album == nil && track.ImageURL != "" {
			storage.FetchImage(util.TRACKS, track.ImageURL, track.FullTitle)
		}
	}

	return nil
}
