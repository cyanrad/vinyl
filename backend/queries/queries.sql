-- name: GetAllTrackItems :many
SELECT  t.id                                       AS track_id,
        t.title,
        al.id                                      AS album_id,
        al.name                                    AS album_name,
        CAST(JSON_GROUP_ARRAY(ar.id) AS TEXT)      AS artist_ids,
        CAST(JSON_GROUP_ARRAY(ar.name) AS TEXT)    AS artist_names
  FROM  tracks AS t
  JOIN  tracks_artists AS tar         ON t.id = tar.track_id
  JOIN  artists AS ar                 ON tar.artist_id = ar.id
  LEFT  JOIN tracks_albums AS tal     ON t.id = tal.track_id
  LEFT  JOIN albums AS al             ON tal.album_id = al.id
 GROUP  BY 1,2,3,4
 ORDER  BY t.created_at;

-- name: GetTrackItemById :one
SELECT  t.id                                       AS track_id,
        t.title,
        al.id                                      AS album_id,
        al.name                                    AS album_name,
        CAST(JSON_GROUP_ARRAY(ar.id) AS TEXT)      AS artist_ids,
        CAST(JSON_GROUP_ARRAY(ar.name) AS TEXT)    AS artist_names
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

-- name: GetAlbumByName :one
SELECT  a.id
  FROM  albums a
 WHERE  a.name = ?;

-- name: GetAlbumById :one
SELECT  al.id,
        al.name,
        al.description,
        al.created_at,
        GROUP_CONCAT(ar.id, ',')     AS artist_ids,
        GROUP_CONCAT(ar.name, ',')   AS artist_names
  FROM  albums al
  LEFT  JOIN artists_albums aa  ON al.id = aa.album_id
  LEFT  JOIN artists ar         ON aa.artist_id = ar.id
 WHERE  al.id = ?
 GROUP  BY 1,2,3,4;
