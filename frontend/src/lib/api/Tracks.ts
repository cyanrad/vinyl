import PocketBase from "pocketbase";

import { ALBUMS_COLLECTION, ARTISTS_COLLECTION, TRACKS_COLLECTION } from "./consts";
import type { Track } from "./Types";
import { generateFileUrl } from "./Util";

async function getAllTracks(pb: PocketBase): Promise<Track[]> {
    const records = await pb.collection(TRACKS_COLLECTION).getFullList<Track>();
    return records;
}

function generateTrackCoverUrl(track: Track): string | null {
    return track.cover ? generateFileUrl(TRACKS_COLLECTION, track.id, track.cover) : null;
}

function generateTrackAudioUrl(track: Track): string | null {
    return generateFileUrl(TRACKS_COLLECTION, track.id, track.audio);
}

export { getAllTracks, generateTrackCoverUrl, generateTrackAudioUrl };
