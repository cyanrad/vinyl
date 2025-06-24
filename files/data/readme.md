# Data Files

In this directory you can add the data files for your artists, tracks, albums, etc.

## Artists

`artists.json`: Contains a list of artists with their details.
`name` of the artist
`description(optional)` of who the artist is and what genre they work in
`links(optional)` currently spotify, soundcloud, & personal are implemented

```json
[
  {
    "name": "example artist",
    "description": "This is an example artist description.",
    "links": {
      "soundcloud": "https://soundcloud.com/example",
      "spotify": "https://open.spotify.com/artist/example",
      "personal": "https://example.com"
    }
  }
]
```

## Tracks

`tracks.json`: Contains a list of tracks with their details.
`title` of the track
`artist(list)` list of the artists who worked on the track. **Must be present in `artists.json`**
`tags(optional list)` to categorise the track

```json
[
  {
    "title": "Self Obsessed",
    "artists": ["Oco."],
    "tags": ["D&B", "breakcore", "electronic"]
  }
]
```
