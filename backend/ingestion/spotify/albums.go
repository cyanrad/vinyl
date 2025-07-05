package spotify

import (
	"log"
	"main/ingestion/storage"
	"main/util"

	"github.com/zmb3/spotify/v2"
)

func (s *SpotifyConn) GetFullAlbums(albumIDs []spotify.ID) ([]*spotify.FullAlbum, error) {
	albumIDs = deduplicate(albumIDs)
	log.Printf("Getting %d albums\n", len(albumIDs))

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
	log.Printf("Generating %d Spotify albums data from public API\n", len(nonCachedIDs))
	for offset := 0; offset < len(nonCachedIDs); offset += util.ALBUM_PAGE_SIZE {
		util.LogProgress(offset, len(nonCachedIDs))
		end := min(offset+util.ALBUM_PAGE_SIZE, len(nonCachedIDs))

		albumsPage, err := s.client.GetAlbums(s.ctx, nonCachedIDs[offset:end])
		if err != nil {
			return nil, err
		}
		albums = append(albums, albumsPage...)
	}

	log.Printf("Spotify albums API complete, Caching %d objects\n", len(nonCachedIDs))
	for i, a := range albums[uncachedStartIndex:] {
		util.LogProgress(i, len(albums)-uncachedStartIndex)
		err := s.cache.Store(util.ALBUMS, util.SOURCE_SPOTIFY, a.ID.String(), a)
		if err != nil {
			return nil, nil
		}
	}

	return albums, nil
}

func generateAlbumIngestion(album spotify.FullAlbum) storage.AlbumIngestion {
	// TODO: add artists logic here
	return storage.AlbumIngestion{
		Name:        album.Name,
		Description: nil,
	}
}
