export interface Track {
    readonly id: string;
    readonly created: string;
    description?: string;
    title: string;
    tags?: Record<string, unknown>; // JSON field
}

export interface Album {
    readonly id: string;
    readonly created_at: string;
    name: string;
    description?: string;
}

export interface Artist {
    readonly id: string;
    readonly created_at: string;
    name: string;
    description?: string;
    links: object;
}

export interface TrackItem {
    // track
    title: string;
    description?: string;
    tags?: Record<string, unknown>; // JSON field

    // album
    albumTitle: string;

    // artist
    artistNames: string[];

    // ids
    trackId: string;
    albumId: string;
    artistIds: string[];
}
