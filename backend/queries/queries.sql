-- name: GetAllTrackItems :many
SELECT  t.id                         AS track_id,
        t.title,
        al.id                        AS album_id,
        al.name                      AS album_name,
        GROUP_CONCAT(ar.id, ',')     AS artist_ids,
        GROUP_CONCAT(ar.name, ',')   AS artist_names
  FROM  tracks AS t
  JOIN  tracks_artists AS tar         ON t.id = tar.track_id
  JOIN  artists AS ar                 ON tar.artist_id = ar.id
  LEFT  JOIN tracks_albums AS tal     ON t.id = tal.track_id
  LEFT  JOIN albums AS al             ON tal.album_id = al.id
 GROUP  BY 1,2,3,4
 ORDER  BY t.created_at;

-- name: GetTrackItemById :one
SELECT  t.id                         AS track_id,
        t.title,
        al.id                        AS album_id,
        al.name                      AS album_name,
        GROUP_CONCAT(ar.id, ',')     AS artist_ids,
        GROUP_CONCAT(ar.name, ',')   AS artist_names
  FROM  tracks AS t
  JOIN  tracks_artists AS tar         ON t.id = tar.track_id
  JOIN  artists AS ar                 ON tar.artist_id = ar.id
  LEFT  JOIN tracks_albums AS tal     ON t.id = tal.track_id
  LEFT  JOIN albums AS al             ON tal.album_id = al.id
 WHERE  t.id = ?
 GROUP  BY 1,2,3,4;

-- name: GetArtistByName :one
SELECT  a.id
  FROM  artists AS a
 WHERE  a.name = ?;

-- name: GetArtistById :one
SELECT  a.*
  FROM  artists a
 WHERE  a.id = ?;

-- name: GetAlbumById :one
SELECT  a.*
  FROM  albums a
 WHERE  a.id = ?;
