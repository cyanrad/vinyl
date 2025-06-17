import { ALBUMS_COLLECTION, ARTISTS_COLLECTION, TRACKS_COLLECTION } from "./consts";
import type { Artist, TrackItem } from "./Types";
import { generateFileUrl } from "./Util";
import type PocketBase from "pocketbase";

export async function getAllTrackItem(pb: PocketBase) {
    const items = await pb.collection(TRACKS_COLLECTION).getFullList({
        expand: "album,artists",
    });

    const trackItems: TrackItem[] = items.map((track) => {
        const artists = track.expand?.artists || [];
        const album = track.expand?.album;

        return {
            // track
            title: track.title,
            description: track.description,
            audio: track.audio,
            cover: track.cover,
            tags: track.tags,

            // album
            albumTitle: album?.title,
            albumCover: album?.cover,

            // artists
            artistNames: artists.map((artist: Artist) => artist.name),
            artistImages: artists.map((artist: Artist) => artist.image),

            // ids
            trackId: track.id,
            albumId: album?.id,
            artistIds: artists.map((artist: Artist) => artist.id),
        };
    });

    return trackItems;
}

export function generateTrackItemCoverUrl(trackItem: TrackItem): string | null {
    return trackItem.cover ? generateFileUrl(TRACKS_COLLECTION, trackItem.trackId, trackItem.cover) : null;
}

export function generateTrackItemAudioUrl(trackItem: TrackItem): string {
    return generateFileUrl(TRACKS_COLLECTION, trackItem.trackId, trackItem.audio);
}

export function generateTrackItemAlbumCoverUrl(trackItem: TrackItem): string | null {
    return trackItem.albumCover ? generateFileUrl(ALBUMS_COLLECTION, trackItem.albumId, trackItem.albumCover) : null;
}

export function generateTrackItemArtistImageUrl(trackItem: TrackItem, index: number): string | null {
    return trackItem.artistImages[index]
        ? generateFileUrl(ARTISTS_COLLECTION, trackItem.artistIds[index], trackItem.artistImages[index])
        : null;
}
