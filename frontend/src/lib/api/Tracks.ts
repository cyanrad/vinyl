import PocketBase from "pocketbase";

import { ALBUMS_COLLECTION, ARTISTS_COLLECTION, TRACKS_COLLECTION } from "./consts";
import type { Track } from "./Types";
import { generateFileUrl } from "./Util";
import { getAlbumById } from "./Albums";
import { getArtistById } from "./Artists";

async function getAllTracks(pb: PocketBase): Promise<Track[]> {
    const records = await pb.collection(TRACKS_COLLECTION).getFullList<Track>();
    return records;
}

async function generateTrackCoverUrl(pb: PocketBase, track: Track): Promise<string | null> {
    let filename: string | null = null;
    let id: string | null = null;
    let collection: string | null = null;

    if (track.cover) {
        filename = track.cover;
        id = track.id;
        collection = TRACKS_COLLECTION;
    }

    // if no track cover, use the album cover
    if (filename === null && track.album) {
        const album = await getAlbumById(pb, track.album);
        if (album.cover) {
            filename = album.cover;
            id = album.id;
            collection = ALBUMS_COLLECTION;
        }
    }

    // if no track & album cover, use the artist image
    if (filename === null) {
        const artist = await getArtistById(pb, track.artist);
        if (artist.image) {
            filename = artist.image;
            id = artist.id;
            collection = ARTISTS_COLLECTION;
        }
    }

    return filename && id && collection ? generateFileUrl(collection, id, filename) : null;
}

function generateTrackAudioUrl(track: Track): string | null {
    return generateFileUrl(TRACKS_COLLECTION, track.id, track.audio);
}

export { getAllTracks, generateTrackCoverUrl, generateTrackAudioUrl };
