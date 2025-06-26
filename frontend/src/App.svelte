<script lang="ts">
    import { onMount, setContext, untrack } from "svelte";

    // styles
    import "./app.css";

    // components
    import Base from "./lib/base/Base.svelte";
    import SideBar from "./lib/side-bar/SideBar.svelte";

    // api
    import { API_URL } from "./lib/api/Consts";
    import { getAllTrackItem, generateTrackItemAudioUrl } from "./lib/api/TrackItems";
    import type { TrackItem } from "./lib/api/Types";
    import PocketBase from "pocketbase";
    import PlayerState from "./lib/PlayerState";

    // Setting the FPS context to be used in animations
    setContext("fps", 100);

    // API instance
    const pb = new PocketBase(API_URL);
    setContext("pb", pb);
    let tracks: TrackItem[] = $state([]);

    // the overall state of the player coordinated with all components
    let playerState: PlayerState = $state(PlayerState.Paused);

    // active track
    let activeTrackIndex: number | null = $state(null);
    let nextTrackSignal: number = $state(0);
    let activeTrack: TrackItem | null = $derived(activeTrackIndex !== null ? tracks[activeTrackIndex] : null);

    // active audio
    let audio: HTMLAudioElement | null = $state(null);
    let currentTime: number = $state(0);
    let duration: number = $state(0);
    let currTimeUpdated: boolean = $state(false);

    onMount(async () => {
        tracks = await getAllTrackItem();
    });

    // next track signal indicates a next/prev track button click
    $effect(() => {
        if (audio === null || activeTrack === null || activeTrackIndex === null) return;

        let newIndex = activeTrackIndex + nextTrackSignal;
        if (newIndex >= 0 && newIndex < tracks.length) {
            activeTrackIndex = newIndex;
        }

        untrack(() => {
            nextTrackSignal = 0;
        });
    });

    // go to the next track when the current track ends
    $effect(() => {
        if (
            audio &&
            duration != 0 && // need to wait for metadata to be loaded even after audio is loaded
            currentTime >= duration &&
            activeTrackIndex !== null &&
            activeTrackIndex < tracks.length - 1
        ) {
            activeTrackIndex++;
        }
    });

    // update player state on active track change
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

        // TODO: this is not enough for generation of audio, search of a better check
        audio = new Audio(generateTrackItemAudioUrl(activeTrack));

        playerState = PlayerState.Playing;
    });

    // load duration when metadata is ready
    $effect(() => {
        if (!audio) return;

        audio.addEventListener("loadedmetadata", () => {
            if (!audio) return;
            duration = audio.duration;
        });

        // update current time as it plays
        audio.addEventListener("timeupdate", () => {
            if (!audio) return;
            currentTime = audio.currentTime;
        });
    });
</script>

<main>
    <SideBar {tracks} bind:activeTrackIndex />
    <div class="flex items-center justify-center h-screen w-screen fixed inset-0">
        <Base
            {activeTrack}
            {duration}
            bind:nextTrackSignal
            bind:playerState
            bind:audio
            bind:currentTime
            bind:currTimeUpdated
        />
    </div>
</main>

<style>
</style>
