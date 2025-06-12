import PocketBase from "pocketbase";

import { ARTISTS_COLLECTION } from "./consts";
import type { Artist } from "./Types";
import { generateFileUrl } from "./Util";

async function getArtistById(pb: PocketBase, artistId: string): Promise<Artist> {
    return await pb.collection(ARTISTS_COLLECTION).getOne<Artist>(artistId);
}

function generateArtistImageUrl(artist: Artist): string | null {
    return artist.image ? generateFileUrl(ARTISTS_COLLECTION, artist.id, artist.image) : null;
}

export { getArtistById, generateArtistImageUrl };
