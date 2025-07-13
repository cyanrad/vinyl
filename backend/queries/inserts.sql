-- name: CreateTrack :one
INSERT INTO tracks (title, full_title, description, tags) VALUES (?, ?, ?, ?) RETURNING *;

-- name: CreateArtist :exec
INSERT INTO artists (name, description, links) VALUES (?, ?, ?);

-- name: CreateAlbum :one
INSERT INTO albums (name, full_name, description) VALUES (?, ?, ?) RETURNING *;

-- name: CreatePlaylist :exec
INSERT INTO playlists (name, description) VALUES (?, ?);

-- name: CreateTrackArtist :exec
INSERT INTO tracks_artists (track_id, artist_id, rank) VALUES (?, ?, ?);

-- name: CreateTrackAlbum :exec
INSERT INTO tracks_albums (track_id, album_id, rank) VALUES (?, ?, ?);

-- name: CreateArtistAlbum :exec
INSERT INTO artists_albums (artist_id, album_id, rank) VALUES (?, ?, ?);

-- name: CreateTrackPlaylist :exec
INSERT INTO tracks_playlists (track_id, playlist_id, rank) VALUES (?, ?, ?);
