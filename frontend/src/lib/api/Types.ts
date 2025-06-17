export interface Track {
    readonly id: string;
    readonly collectionId: string;
    readonly collectionName: string;
    readonly created: string;
    readonly updated: string;
    artists: string[]; // relation IDs
    audio: string; // filename
    album?: string;
    cover?: string; // filename
    description?: string;
    title: string;
    tags?: Record<string, unknown>; // JSON field
}

export interface Album {
    readonly id: string;
    readonly collectionId: string;
    readonly collectionName: string;
    readonly created: string;
    readonly updated: string;
    artist: string; // relation ID
    cover?: string; // filename
    description?: string;
    title: string;
    tags?: Record<string, unknown>; // JSON field
}

export interface Artist {
    readonly id: string;
    readonly collectionId: string;
    readonly collectionName: string;
    readonly created: string;
    readonly updated: string;
    image?: string; // filename
    description?: string;
    name: string;
    links?: Record<string, unknown>; // JSON field
}

export interface TrackItem {
    // track
    title: string;
    description?: string;
    audio: string; // filename
    cover?: string; // filename
    tags?: Record<string, unknown>; // JSON field

    // album
    albumTitle: string;
    albumCover?: string; // filename

    // artist
    artistNames: string[];
    artistImages: string[]; // filenames

    // ids
    trackId: string;
    albumId: string;
    artistIds: string[];
}
