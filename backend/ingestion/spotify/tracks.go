package spotify

import (
	"log"
	"main/util"

	"github.com/zmb3/spotify/v2"
)

// TODO: this function is copy pasted word for word 3 times, generalize it
func (s *SpotifyConn) GetFullTracks(trackIDs []spotify.ID) ([]*spotify.FullTrack, error) {
	trackIDs = deduplicate(trackIDs)
	log.Printf("Getting %d tracks\n", len(trackIDs))

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

	// return if all tracks are cached
	if len(trackIDs) == len(tracks) {
		return tracks, nil
	}

	// getting all tracks from an external API
	uncachedStartIndex := len(tracks)
	log.Printf("Generating %d Spotify tracks data from public API\n", len(nonCachedIDs))
	for offset := 0; offset < len(nonCachedIDs); offset += util.TRACK_PAGE_SIZE {
		util.LogProgress(offset, len(nonCachedIDs))
		end := min(offset+util.TRACK_PAGE_SIZE, len(nonCachedIDs))

		tracksPage, err := s.client.GetTracks(s.ctx, nonCachedIDs[offset:end])
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, tracksPage...)
	}

	log.Printf("Spotify tracks API complete, Caching %d objects\n", len(nonCachedIDs))
	for i, a := range tracks[uncachedStartIndex:] {
		util.LogProgress(i, len(tracks)-uncachedStartIndex)
		err := s.cache.Store(util.TRACKS, util.SOURCE_SPOTIFY, a.ID.String(), a)
		if err != nil {
			return nil, nil
		}
	}

	return tracks, nil
}
