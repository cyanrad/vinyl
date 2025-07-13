package spotify

import (
	"main/ingestion/storage"
	"main/util"

	"github.com/zmb3/spotify/v2"
)

func (s *SpotifyConn) GetFullArtists(artistIDs []spotify.ID) ([]*spotify.FullArtist, error) {
	util.Log.Debugf("Artists count before deduping: %d", len(artistIDs))
	artistIDs = deduplicate(artistIDs)
	util.Log.Infof("Getting %d artists\n", len(artistIDs))

	// retrun array
	artists := make([]*spotify.FullArtist, 0, len(artistIDs))
	// to store the IDs that were not cached so we can fetch from server
	nonCachedIDs := make([]spotify.ID, 0, len(artistIDs)) // not really the best idea, but keeps the code simple

	// getting all cached artists
	for _, id := range artistIDs {
		data := spotify.FullArtist{}
		found, err := s.getCached(util.ARTISTS, util.SOURCE_SPOTIFY, id.String(), &data)
		if err != nil {
			return nil, err
		} else if found {
			artists = append(artists, &data)
		} else {
			nonCachedIDs = append(nonCachedIDs, id)
		}
	}

	// return if all artists are cached
	if len(artistIDs) == len(artists) {
		return artists, nil
	}

	// getting all artists from an external API
	uncachedStartIndex := len(artists)
	util.Log.Infof("Generating %d Spotify artists data from public API\n", len(nonCachedIDs))
	for offset := 0; offset < len(nonCachedIDs); offset += util.ARTIST_PAGE_SIZE {
		util.LogProgress(offset, len(nonCachedIDs))
		end := min(offset+util.ARTIST_PAGE_SIZE, len(nonCachedIDs))

		artistsPage, err := s.client.GetArtists(s.ctx, nonCachedIDs[offset:end]...)
		if err != nil {
			return nil, err
		}
		artists = append(artists, artistsPage...)
	}

	util.Log.Infof("Spotify artists API complete, Caching %d objects\n", len(nonCachedIDs))
	for i, a := range artists[uncachedStartIndex:] {
		util.LogProgress(i, len(artists)-uncachedStartIndex)
		err := s.cache.Store(util.ARTISTS, util.SOURCE_SPOTIFY, a.ID.String(), a)
		if err != nil {
			return nil, nil
		}
	}

	return artists, nil
}

func (s *SpotifyConn) GenerateArtistSpotifyIngestion(artists []*spotify.FullArtist) SpotifyIngestion {
	return SpotifyIngestion{
		Artists: GenerateArtistIngestions(artists),
		Albums:  nil,
		Tracks:  nil,
	}
}

func (s *SpotifyConn) SimpleToFullArtists(artists []spotify.SimpleArtist) ([]*spotify.FullArtist, error) {
	// getting artist IDs
	artistIDs := make([]spotify.ID, len(artists))
	for i, a := range artists {
		artistIDs[i] = a.ID
	}

	return s.GetFullArtists(artistIDs)
}

func GenerateArtistIngestions(artists []*spotify.FullArtist) []ArtistIngestion {
	// creating track artist ingestion
	ingestions := make([]ArtistIngestion, len(artists))
	for i, a := range artists {
		links := storage.ArtistLinks{}
		if url, ok := a.ExternalURLs["spotify"]; ok {
			links.Spotify = &url
		}

		imageURL := ""
		if len(a.Images) > 0 {
			imageURL = a.Images[0].URL
		}
		ingestions[i] = ArtistIngestion{
			Name:     util.GenerateArtistName(a.Name),
			Links:    links,
			ImageURL: imageURL,
		}
	}

	return ingestions
}
