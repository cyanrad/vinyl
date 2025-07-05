package spotify

import (
	"log"
	"main/util"

	"github.com/zmb3/spotify/v2"
)

type PlaylistData struct {
	ArtistIDs []spotify.ID
	AlbumIDs  []spotify.ID
	Tracks    []spotify.FullTrack
}

func (s *SpotifyConn) GeneratePlaylistData(playlistID string) (PlaylistData, error) {
	data := PlaylistData{}
	found, err := s.getCached(util.PLAYLISTS, util.SOURCE_SPOTIFY, playlistID, &data)
	if err != nil {
		return PlaylistData{}, err
	} else if found {
		return data, nil
	}

	log.Printf("Generating Spotify playlist %s data from public API\n", playlistID)
	playlistPage, err := s.client.GetPlaylist(s.ctx, spotify.ID(playlistID))
	if err != nil {
		return PlaylistData{}, err
	}

	offset := util.PLAYLIST_PAGE_SIZE
	playlistSize := int(playlistPage.Tracks.Total)

	tracks := make([]spotify.PlaylistTrack, 0, playlistSize)
	tracks = append(tracks, playlistPage.Tracks.Tracks...)

	for ; offset < playlistSize; offset += util.PLAYLIST_PAGE_SIZE {
		util.LogProgress(offset, playlistSize)
		playlistPage, err = s.client.GetPlaylist(s.ctx, spotify.ID(playlistID), spotify.Offset(offset))
		if err != nil {
			return PlaylistData{}, err
		}
		tracks = append(tracks, playlistPage.Tracks.Tracks...)
	}

	log.Println("Spotify playlist API calls complete, generating data")
	playlistData := PlaylistData{
		ArtistIDs: []spotify.ID{}, // we can't pre allocat as total count is unknown
		AlbumIDs:  make([]spotify.ID, len(tracks)),
		Tracks:    make([]spotify.FullTrack, len(tracks)),
	}
	for i, t := range tracks {
		playlistData.Tracks[i] = t.Track
		playlistData.AlbumIDs[i] = t.Track.Album.ID
		for _, a := range t.Track.Artists {
			playlistData.ArtistIDs = append(playlistData.ArtistIDs, a.ID)
		}
	}

	log.Printf("Spotify playlist %s complete. Caching object\n", playlistID)
	err = s.cache.Store(util.PLAYLISTS, util.SOURCE_SPOTIFY, playlistID, playlistData)
	if err != nil {
		return PlaylistData{}, nil
	}

	return playlistData, nil
}
