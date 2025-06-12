import PocketBase from "pocketbase";

import { ALBUMS_COLLECTION } from "./consts";
import type { Album } from "./Types";

async function getAlbumById(pb: PocketBase, albumId: string): Promise<Album> {
    return await pb.collection(ALBUMS_COLLECTION).getOne<Album>(albumId);
}

export { getAlbumById };
