<script lang="ts">
    import { onMount, setContext, untrack } from "svelte";

    // styles
    import "./app.css";

    // components
    import Base from "./lib/Base.svelte";
    import SideBar from "./lib/SideBar.svelte";

    // api
    import { API_URL } from "./lib/api/consts";
    import { getAllTracks, generateTrackAudioUrl } from "./lib/api/Tracks";
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

    let activeTrackIndex: number | null = $state(null);
    let activeTrack: Track | null = $state(null);
    let activeArtist: Artist | null = $state(null);
    let activeAlbum: Album | null = $state(null);

    let audio: HTMLAudioElement | null = $state(null);

    onMount(async () => {
        tracks = await getAllTracks(pb);
    });

    $effect(() => {
        if (!activeTrackIndex) return;

        untrack(() => {
            if (audio) {
                audio.pause(); // Stop the audio
                audio.src = ""; // Clear the source to stop downloading
                audio.load(); // Reset the audio element
                audio = null; // Remove reference to allow GC
            }
        });

        activeTrack = tracks[activeTrackIndex];
        getArtistById(pb, activeTrack.artist).then((artist) => {
            activeArtist = artist;
        });

        if (activeTrack.album) {
            getAlbumById(pb, activeTrack.album).then((album) => {
                activeAlbum = album;
            });
        }

        audio = new Audio(generateTrackAudioUrl(activeTrack));
    });
</script>

<main>
    <SideBar {tracks} bind:activeTrackIndex />
    <div class="flex items-center justify-center h-screen w-screen fixed inset-0">
        <Base {activeTrack} {activeArtist} {activeAlbum} bind:audio />
    </div>
</main>

<style>
</style>
