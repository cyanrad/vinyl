package spotify

import (
	"main/util"

	"github.com/zmb3/spotify/v2"
)

func (s *SpotifyConn) GetFullAlbums(albumIDs []spotify.ID) ([]*spotify.FullAlbum, error) {
	util.Log.Debugf("Albums count before deduping: %d", len(albumIDs))
	albumIDs = deduplicate(albumIDs)
	util.Log.Infof("Getting %d albums\n", len(albumIDs))

	// retrun array
	albums := make([]*spotify.FullAlbum, 0, len(albumIDs))
	// to store the IDs that were not cached so we can fetch from server
	nonCachedIDs := make([]spotify.ID, 0, len(albumIDs)) // not really the best idea, but keeps the code simple

	// getting all cached albums
	for _, id := range albumIDs {
		data := spotify.FullAlbum{}
		found, err := s.getCached(util.ALBUMS, util.SOURCE_SPOTIFY, id.String(), &data)
		if err != nil {
			return nil, err
		} else if found {
			albums = append(albums, &data)
		} else {
			nonCachedIDs = append(nonCachedIDs, id)
		}
	}

	// return if all albums are cached
	if len(albumIDs) == len(albums) {
		return albums, nil
	}

	// getting all albums from an external API
	uncachedStartIndex := len(albums)
	util.Log.Infof("Generating %d Spotify albums data from public API\n", len(nonCachedIDs))
	for offset := 0; offset < len(nonCachedIDs); offset += util.ALBUM_PAGE_SIZE {
		util.LogProgress(offset, len(nonCachedIDs))
		end := min(offset+util.ALBUM_PAGE_SIZE, len(nonCachedIDs))

		albumsPage, err := s.client.GetAlbums(s.ctx, nonCachedIDs[offset:end])
		if err != nil {
			return nil, err
		}
		albums = append(albums, albumsPage...)
	}

	util.Log.Infof("Spotify albums API complete, Caching %d objects\n", len(nonCachedIDs))
	for i, a := range albums[uncachedStartIndex:] {
		util.LogProgress(i, len(albums)-uncachedStartIndex)
		err := s.cache.Store(util.ALBUMS, util.SOURCE_SPOTIFY, a.ID.String(), a)
		if err != nil {
			return nil, nil
		}
	}

	return albums, nil
}

func (s *SpotifyConn) GenerateAlbumSpotifyIngestion(albums []*spotify.FullAlbum) (SpotifyIngestion, error) {
	// generating this list as to save on external request count
	simpleArtists := []spotify.SimpleArtist{}
	fullAlbums := make([]*spotify.FullAlbum, 0, len(albums))
	for _, album := range albums {
		simpleArtists = append(simpleArtists, album.Artists...)

		// if the album is a single we don't count it as an album in our data
		if album.AlbumType != "single" {
			fullAlbums = append(fullAlbums, album)
		}
	}

	fullArtists, err := s.SimpleToFullArtists(simpleArtists)
	if err != nil {
		return SpotifyIngestion{}, err
	}

	return SpotifyIngestion{
		Artists: GenerateArtistIngestions(fullArtists),
		Albums:  GenerateAlbumIngestion(fullAlbums),
		Tracks:  nil,
	}, nil
}

func (s *SpotifyConn) SimpleToFullAlbums(albums []spotify.SimpleAlbum) ([]*spotify.FullAlbum, error) {
	// getting album IDs
	albumIDs := make([]spotify.ID, len(albums))
	for i, a := range albums {
		albumIDs[i] = a.ID
	}

	return s.GetFullAlbums(albumIDs)
}

func GenerateAlbumIngestion(albums []*spotify.FullAlbum) []AlbumIngestion {
	// looping and creating ingestions
	ingestions := make([]AlbumIngestion, 0, len(albums))
	for _, album := range albums {
		// if the album is a single we don't count it as an album in our data
		if album.AlbumType == "single" {
			continue
		}

		// getting related artists
		artistNames := make([]string, len(album.Artists))
		for i, artist := range album.Artists {
			artistNames[i] = util.GenerateArtistName(artist.Name)
		}

		// generating album ingestion
		imageURL := ""
		if len(album.Images) > 0 {
			imageURL = album.Images[0].URL
		}
		albumIngestion := AlbumIngestion{
			Name:     util.CleanName(album.Name),
			FullName: util.GenerateAlbumName(artistNames, album.Name),
			ImageURL: imageURL,
			Artists:  artistNames,
		}

		ingestions = append(ingestions, albumIngestion)
	}

	return ingestions
}
