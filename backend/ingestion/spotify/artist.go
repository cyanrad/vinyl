package spotify

import (
	"log"
	"main/ingestion/storage"
	"main/util"

	"github.com/zmb3/spotify/v2"
)

func (s *SpotifyConn) GetFullArtists(artistIDs []spotify.ID) ([]*spotify.FullArtist, error) {
	artistIDs = deduplicate(artistIDs)
	log.Printf("Getting %d artists\n", len(artistIDs))

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
	log.Printf("Generating %d Spotify artists data from public API\n", len(nonCachedIDs))
	for offset := 0; offset < len(nonCachedIDs); offset += util.ARTIST_PAGE_SIZE {
		util.LogProgress(offset, len(nonCachedIDs))
		end := min(offset+50, len(nonCachedIDs))

		artistsPage, err := s.client.GetArtists(s.ctx, nonCachedIDs[offset:end]...)
		if err != nil {
			return nil, err
		}
		artists = append(artists, artistsPage...)
	}

	log.Printf("Spotify artists API complete, Caching %d objects\n", len(nonCachedIDs))
	for i, a := range artists[uncachedStartIndex:] {
		util.LogProgress(i, len(artists)-uncachedStartIndex)
		err := s.cache.Store(util.ARTISTS, util.SOURCE_SPOTIFY, a.ID.String(), a)
		if err != nil {
			return nil, nil
		}
	}

	return artists, nil
}

func generateArtistIngestion(artist spotify.FullArtist) storage.ArtistIngestion {
	spotifyURL := util.SPOTIFY_ARTIST_URL_BASE + artist.ID.String()

	return storage.ArtistIngestion{
		Name:        artist.Name,
		Description: nil,
		Links: storage.ArtistLinks{
			Spotify: &spotifyURL,
		},
	}
}

// WARNING: please for the love of god find a better way to do this
func deduplicate(input []spotify.ID) []spotify.ID {
	seen := make(map[spotify.ID]bool)
	result := []spotify.ID{}

	for _, val := range input {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}

	return result
}
