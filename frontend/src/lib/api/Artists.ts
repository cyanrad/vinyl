import PocketBase from "pocketbase";

import { ARTISTS_COLLECTION } from "./consts";
import type { Artist } from "./Types";

async function getArtistById(pb: PocketBase, artistId: string): Promise<Artist> {
    return await pb.collection(ARTISTS_COLLECTION).getOne<Artist>(artistId);
}

export { getArtistById };
