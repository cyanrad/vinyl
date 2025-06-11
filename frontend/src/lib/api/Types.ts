export interface Track {
    readonly id: string;
    readonly collectionId: string;
    readonly collectionName: string;
    readonly created: string;
    readonly updated: string;
    artist: string; // relation ID
    audio: string; // filename
    album?: string;
    cover?: string; // filename
    description?: string;
    title: string;
    tags?: Record<string, unknown>; // JSON field
}
