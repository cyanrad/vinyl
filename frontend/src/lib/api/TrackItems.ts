import type { TrackItem } from "./Types";
import { TRACK_ITEMS_RESOURCE, API_URL, ALBUMS_RESOURCE, ARTISTS_RESOURCE, TRACKS_RESOURCE } from "./Consts";
import { generateFileUrl } from "./Util";

export async function getAllTrackItem(): Promise<TrackItem[]> {
    const response = await fetch(API_URL + "/" + TRACK_ITEMS_RESOURCE);

    if (!response.ok) {
        throw new Error(`fetching track items failed: ${response.status}`);
    }

    const items = await response.json();

    // converting artsitIds & Names from string of an array to JSON array
    const processedData = items.map((item: any) => ({
        ...item,
        artistNames: typeof item.artistNames === "string" ? JSON.parse(item.artistNames) : item.artistNames,
        artistIds: typeof item.artistIds === "string" ? JSON.parse(item.artistIds) : item.artistIds,
    })) as TrackItem[];

    return processedData;
}

export function generateTrackItemAlbumCoverUrl(item: TrackItem): string | null {
    if (item.albumId === null) return null;

    return generateFileUrl(ALBUMS_RESOURCE, item.albumId);
}

export function generateTrackItemArtistImageUrl(item: TrackItem): string {
    return generateFileUrl(ARTISTS_RESOURCE, item.artistIds[0]);
}

export function generateTrackItemCoverUrl(item: TrackItem): string {
    return generateFileUrl(TRACKS_RESOURCE, item.trackId);
}

export function generateTrackItemAudioUrl(item: TrackItem): string {
    return `${API_URL}/${TRACKS_RESOURCE}/${item.trackId}/audio`;
}
