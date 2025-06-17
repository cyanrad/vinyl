<script lang="ts">
    import { onMount, setContext, untrack } from "svelte";

    // styles
    import "./app.css";

    // components
    import Base from "./lib/Base.svelte";
    import SideBar from "./lib/SideBar.svelte";

    // api
    import { API_URL } from "./lib/api/consts";
    import { getAllTrackItem, generateTrackItemAudioUrl } from "./lib/api/TrackItems";
    import type { TrackItem } from "./lib/api/Types";
    import PocketBase from "pocketbase";

    // Setting the FPS context to be used in animations
    setContext("fps", 100);

    // API instance
    const pb = new PocketBase(API_URL);
    setContext("pb", pb);
    let tracks: TrackItem[] = $state([]);

    let activeTrackIndex: number | null = $state(null);
    let activeTrack: TrackItem | null = $derived(activeTrackIndex !== null ? tracks[activeTrackIndex] : null);

    let audio: HTMLAudioElement | null = $state(null);
    let currentTime: number = $state(0);
    let duration: number = $state(0);
    let currTimeUpdated: boolean = $state(false);

    onMount(async () => {
        tracks = await getAllTrackItem(pb);
    });

    $effect(() => {
        if (activeTrackIndex === undefined || activeTrackIndex === null || !activeTrack) return;

        // resetting audio variables
        untrack(() => {
            if (audio) {
                audio.pause(); // Stop the audio
                audio.src = ""; // Clear the source to stop downloading
                audio.load(); // Reset the audio element
                audio = null; // Remove reference to allow GC
            }

            currentTime = 0;
            duration = 0;
            currTimeUpdated = true;
        });

        audio = new Audio(generateTrackItemAudioUrl(activeTrack));
    });
</script>

<main>
    <SideBar {tracks} bind:activeTrackIndex />
    <div class="flex items-center justify-center h-screen w-screen fixed inset-0">
        <Base {activeTrack} {duration} bind:audio bind:currentTime bind:currTimeUpdated />
    </div>
</main>

<style>
</style>
