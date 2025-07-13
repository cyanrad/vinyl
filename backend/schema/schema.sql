CREATE TABLE tracks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    full_title VARCHAR(255) NOT NULL,
    description TEXT,
    tags TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE artists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    links TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE albums (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    full_name VARCHAR(225) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE playlists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE tracks_artists (
    track_id INTEGER NOT NULL,
    artist_id INTEGER NOT NULL,
    rank INTEGER NOT NULL,
    PRIMARY KEY (track_id, artist_id),
    FOREIGN KEY (track_id) REFERENCES tracks(id),
    FOREIGN KEY (artist_id) REFERENCES artists(id)
);

CREATE TABLE tracks_albums (
    track_id INTEGER NOT NULL,
    album_id INTEGER NOT NULL,
    rank INTEGER NOT NULL,
    PRIMARY KEY (track_id, album_id),
    FOREIGN KEY (track_id) REFERENCES tracks(id),
    FOREIGN KEY (album_id) REFERENCES albums(id)
);

CREATE TABLE artists_albums (
    artist_id INTEGER NOT NULL,
    album_id INTEGER NOT NULL,
    rank INTEGER NOT NULL,
    PRIMARY KEY (artist_id, album_id),
    FOREIGN KEY (artist_id) REFERENCES artists(id),
    FOREIGN KEY (album_id) REFERENCES albums(id)
);

CREATE TABLE tracks_playlists (
    track_id INTEGER NOT NULL,
    playlist_id INTEGER NOT NULL,
    rank INTEGER NOT NULL,
    PRIMARY KEY (track_id, playlist_id),
    FOREIGN KEY (track_id) REFERENCES tracks(id),
    FOREIGN KEY (playlist_id) REFERENCES playlists(id)
);
