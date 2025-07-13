package spotify

import (
	"main/util"

	"github.com/zmb3/spotify/v2"
)

func (s *SpotifyConn) GetFullTracks(trackIDs []spotify.ID) ([]*spotify.FullTrack, error) {
	util.Log.Debugf("Track count before deduping: %d", len(trackIDs))
	trackIDs = deduplicate(trackIDs)
	util.Log.Infof("Getting %d tracks\n", len(trackIDs))

	// retrun array
	tracks := make([]*spotify.FullTrack, 0, len(trackIDs))
	// to store the IDs that were not cached so we can fetch from server
	nonCachedIDs := make([]spotify.ID, 0, len(trackIDs)) // not really the best idea, but keeps the code simple

	// getting all cached tracks
	for _, id := range trackIDs {
		data := spotify.FullTrack{}
		found, err := s.getCached(util.TRACKS, util.SOURCE_SPOTIFY, id.String(), &data)
		if err != nil {
			return nil, err
		} else if found {
			tracks = append(tracks, &data)
		} else {
			nonCachedIDs = append(nonCachedIDs, id)
		}
	}
	util.Log.Debugf("Track caching check complete. %d cached out of %d", len(tracks), len(trackIDs))

	// return if all tracks are cached
	if len(trackIDs) == len(tracks) {
		return tracks, nil
	}

	// getting all tracks from an external API
	uncachedStartIndex := len(tracks)
	util.Log.Infof("Generating %d Spotify tracks data from public API\n", len(nonCachedIDs))
	for offset := 0; offset < len(nonCachedIDs); offset += util.TRACK_PAGE_SIZE {
		util.LogProgress(offset, len(nonCachedIDs))
		end := min(offset+util.TRACK_PAGE_SIZE, len(nonCachedIDs))

		tracksPage, err := s.client.GetTracks(s.ctx, nonCachedIDs[offset:end])
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, tracksPage...)
	}

	util.Log.Infof("Spotify tracks API complete, Caching %d objects\n", len(nonCachedIDs))
	for i, a := range tracks[uncachedStartIndex:] {
		util.LogProgress(i, len(tracks)-uncachedStartIndex)
		util.Log.Debugf("caching track '%s'", a.ID.String())

		err := s.cache.Store(util.TRACKS, util.SOURCE_SPOTIFY, a.ID.String(), a)
		if err != nil {
			return nil, nil
		}
	}

	return tracks, nil
}

func (s *SpotifyConn) GenerateTrackSpotifyIngestion(tracks []*spotify.FullTrack) (SpotifyIngestion, error) {
	// generating this list as to save on external request count
	simpleArtists := []spotify.SimpleArtist{}
	simpleAlbums := []spotify.SimpleAlbum{}
	for _, track := range tracks {
		simpleArtists = append(simpleArtists, track.Artists...)
		simpleArtists = append(simpleArtists, track.Album.Artists...)

		// if the album is a single we don't count it as an album in our data
		if track.Album.AlbumType != "single" {
			simpleAlbums = append(simpleAlbums, track.Album)
		}
	}

	fullArtists, err := s.SimpleToFullArtists(simpleArtists)
	if err != nil {
		return SpotifyIngestion{}, err
	}

	fullAlbums, err := s.SimpleToFullAlbums(simpleAlbums)
	if err != nil {
		return SpotifyIngestion{}, err
	}

	return SpotifyIngestion{
		Artists: GenerateArtistIngestions(fullArtists),
		Albums:  GenerateAlbumIngestion(fullAlbums),
		Tracks:  GenerateTrackIngestion(tracks),
	}, nil
}

func GenerateTrackIngestion(tracks []*spotify.FullTrack) []TrackIngestion {
	ingestions := make([]TrackIngestion, len(tracks))

	for i, track := range tracks {
		var album *string = nil
		if track.Album.AlbumType != "single" {
			albumArtists := make([]string, len(track.Album.Artists))
			for i, artist := range track.Album.Artists {
				albumArtists[i] = util.GenerateArtistName(artist.Name)
			}

			temp := util.GenerateAlbumName(albumArtists, track.Album.Name)
			album = &temp
		}

		artists := make([]string, len(track.Artists))
		for i, artist := range track.Artists {
			artists[i] = util.GenerateArtistName(artist.Name)
		}

		imageURL := ""
		if len(track.Album.Images) > 0 {
			imageURL = track.Album.Images[0].URL
		}
		ingestions[i] = TrackIngestion{
			Title:     util.CleanName(track.Name),
			FullTitle: util.GenerateTrackName(artists, util.CleanName(track.Name)),
			Tags:      nil,
			AlbumRank: int(track.TrackNumber),
			ImageURL:  imageURL,
			Artists:   artists,
			Album:     album,
		}
	}

	return ingestions
}
