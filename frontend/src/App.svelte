<script lang="ts">
    import { onMount, setContext } from "svelte";

    // styles
    import "./app.css";

    // components
    import Base from "./lib/Base.svelte";
    import SideBar from "./lib/SideBar.svelte";

    // api
    import { API_URL } from "./lib/api/consts";
    import { getAllTracks } from "./lib/api/Tracks";
    import { getAlbumById } from "./lib/api/Albums";
    import { getArtistById } from "./lib/api/Artists";
    import type { Track, Artist, Album } from "./lib/api/Types";
    import PocketBase from "pocketbase";

    // Setting the FPS context to be used in animations
    setContext("fps", 100);

    // API instance
    const pb = new PocketBase(API_URL);
    setContext("pb", pb);
    let tracks: Track[] = $state([]);
    let activeTrack: Track | null = $state(null);
    let activeArtist: Artist | null = $state(null);
    let activeAlbum: Album | null = $state(null);

    onMount(async () => {
        tracks = await getAllTracks(pb);
    });
</script>

<main>
    <SideBar {tracks} bind:activeTrack bind:activeArtist bind:activeAlbum />
    <div class="flex items-center justify-center h-screen w-screen fixed inset-0">
        <Base {activeTrack} {activeArtist} {activeAlbum} />
    </div>
</main>

<style>
</style>
