package spotify

import (
	"log"
	"main/util"

	"github.com/zmb3/spotify/v2"
)

func (s *SpotifyConn) GetPlaylistData(playlistID string) ([]*spotify.FullTrack, error) {
	data := []*spotify.FullTrack{}

	// checking cache
	found, err := s.getCached(util.PLAYLISTS, util.SOURCE_SPOTIFY, playlistID, &data)
	if err != nil {
		return nil, err
	} else if found {
		return data, nil
	}

	// initial page fetch
	log.Printf("Generating Spotify playlist %s data from public API\n", playlistID)
	playlist, err := s.client.GetPlaylist(s.ctx, spotify.ID(playlistID))
	if err != nil {
		return nil, err
	}

	// initializing vars
	offset := util.PLAYLIST_PAGE_SIZE
	playlistSize := int(playlist.Tracks.Total)
	data = make([]*spotify.FullTrack, 0, playlistSize)

	for _, track := range playlist.Tracks.Tracks {
		data = append(data, &track.Track)
	}

	// looping page fetches
	for ; offset < playlistSize; offset += util.PLAYLIST_PAGE_SIZE {
		util.LogProgress(offset, playlistSize)
		playlistPage, err := s.client.GetPlaylistItems(s.ctx, spotify.ID(playlistID), spotify.Offset(offset))
		if err != nil {
			return nil, err
		}

		for _, track := range playlistPage.Items {
			data = append(data, track.Track.Track)
		}
	}

	// caching object
	log.Printf("Spotify playlist %s complete. Caching object\n", playlistID)
	err = s.cache.Store(util.PLAYLISTS, util.SOURCE_SPOTIFY, playlistID, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
