-- name: CreateTrack :exec
INSERT INTO tracks (title, description, tags) VALUES (?, ?, ?);

-- name: CreateArtist :exec
INSERT INTO artists (name, description, links) VALUES (?, ?, ?);

-- name: CreateAlbum :exec
INSERT INTO albums (name, description) VALUES (?, ?);

-- name: CreatePlaylist :exec
INSERT INTO playlists (name, description) VALUES (?, ?);

-- name: CreateTrackArtist :exec
INSERT INTO tracks_artists (track_id, artist_id) VALUES (?, ?);

-- name: CreateTrackAlbum :exec
INSERT INTO tracks_albums (track_id, album_id) VALUES (?, ?);

-- name: CreateArtistAlbum :exec
INSERT INTO artists_albums (artist_id, album_id) VALUES (?, ?);

-- name: CreateTrackPlaylist :exec
INSERT INTO tracks_playlists (track_id, playlist_id) VALUES (?, ?);