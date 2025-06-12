import PocketBase from "pocketbase";

import { ALBUMS_COLLECTION } from "./consts";
import type { Album } from "./Types";
import { generateFileUrl } from "./Util";

async function getAlbumById(pb: PocketBase, albumId: string): Promise<Album> {
    return await pb.collection(ALBUMS_COLLECTION).getOne<Album>(albumId);
}

function generateAlbumCoverUrl(album: Album): string | null {
    return album.cover ? generateFileUrl(ALBUMS_COLLECTION, album.id, album.cover) : null;
}

export { getAlbumById, generateAlbumCoverUrl };
