<script lang="ts">
    import { onMount, getContext } from "svelte";

    // components
    import Vinyl from "./Vinyl.svelte";
    import Controller from "./Controller.svelte";
    import Scrubber from "./Scrubber.svelte";
    import PlayerHead from "./PlayerHead.svelte";
    import AudioControl from "./AudioControl.svelte";

    // state
    import PlayerState from "./PlayerState";
    import { getMovedAudioTime, getNewAudioTime } from "./Audio";

    // api
    import { getAllTracks, generateTrackCoverUrl, generateTrackAudioUrl } from "./api/Tracks";
    import type { Track } from "./api/Types";
    import PocketBase from "pocketbase";

    // API instance
    const pb: PocketBase = getContext("pb");
    let track: Track | null = $state(null);

    // media
    let trackCover: string | null = $derived(track ? generateTrackCoverUrl(track) : null);
    let trackAudio: string | null = $derived(track ? generateTrackAudioUrl(track) : null);

    // the overall state of the player coordinated with all components
    let playerState: PlayerState = $state(PlayerState.Paused);

    // aduio variables
    let audio: HTMLAudioElement | null = $derived(trackAudio ? new Audio(trackAudio) : null);
    let currentTime: number = $state(0);
    let duration: number = $state(0);
    let volume: number = $state(0.5);

    // detect when a scrubber event happens
    let currTimeUpdated: boolean = $state(false);

    // keyboard events
    onMount(async () => {
        // get the track from the api
        const tracks = await getAllTracks(pb);
        track = tracks[0];

        // play/pause events
        document.addEventListener("keydown", (event) => {
            if (event.code === "Space" || event.code === "KeyK") {
                event.preventDefault(); // Prevent page scroll
                playerState = playerState === PlayerState.Playing ? PlayerState.Paused : PlayerState.Playing;
            }
        });

        // forward/rewind events
        document.addEventListener("keydown", (event) => {
            if (event.code === "ArrowRight") {
                currentTime = getMovedAudioTime(currentTime, duration, 5);
                currTimeUpdated = true;
            } else if (event.code === "ArrowLeft") {
                currentTime = getMovedAudioTime(currentTime, duration, -5);
                currTimeUpdated = true;
            }
        });
    });

    // attaching event listeners to the audio element
    $effect(() => {
        if (!audio) return;
        console.log("triggered");

        // load duration when metadata is ready
        audio.addEventListener("loadedmetadata", () => {
            duration = audio.duration;
        });

        // update current time as it plays
        audio.addEventListener("timeupdate", () => {
            currentTime = audio.currentTime;
        });
    });

    // handle audio playback based on player state
    $effect(() => {
        if (!audio) return;

        // update current time if any change happened from scrubber/keyboard
        if (currTimeUpdated) {
            currentTime = getNewAudioTime(currentTime, duration, currentTime);
            audio.currentTime = currentTime;
            currTimeUpdated = false; // reset the flag after updating
        }

        // play/pause based on player state
        if (playerState === PlayerState.Playing) {
            audio.play();
        } else if (playerState === PlayerState.Paused) {
            audio.pause();
        }

        // stop playback if we reach the end
        if (currentTime && duration && currentTime >= duration) {
            playerState = PlayerState.Paused;
            currentTime = 0;
            currTimeUpdated = true;
        }
    });

    // handle volume change
    $effect(() => {
        if (!audio) return;
        audio.volume = volume;
    });

    // base dimensions
    const baseHeight = 620;

    // controller dimensions
    const controllerTop = 520;
    const controllerLeft = 620;

    // scrubber dimensions
    const scrubberTop = 486;
    const scrubberLeft = 564;

    // player head dimensions
    const playerHeadTop = 50;
    const playerHeadLeft = 480;

    // volume control dimensions
    const volumeControlTop = 70;
    const volumeControlLeft = 762;

    // component images
    const baseImage = "base/base.svg";
</script>

<div class="relative">
    <!-- base -->
    <img
        src={baseImage}
        alt="Base"
        class="w-auto select-none pointer-events-none"
        style="height: {baseHeight}px;"
        draggable="false"
    />

    <!-- player head -->
    <div class="absolute z-30" style="top: {playerHeadTop}px; left: {playerHeadLeft}px;">
        <PlayerHead {currentTime} {duration} {playerState} {currTimeUpdated} />
    </div>

    <!-- vinyl -->
    <Vinyl {playerState} {trackCover} {currentTime} {currTimeUpdated} />

    <!-- controller -->
    <div class="absolute" style="top: {controllerTop}px; left: {controllerLeft}px;">
        <Controller bind:playerState bind:currentTime bind:currTimeUpdated {duration} />
    </div>

    <!-- scrubber -->
    <div class="absolute z-40" style="top: {scrubberTop}px; left: {scrubberLeft}px;">
        <Scrubber bind:currentTime bind:currTimeUpdated {duration} />
    </div>

    <!-- volume control -->
    <div class="absolute" style="top: {volumeControlTop}px; left: {volumeControlLeft}px;">
        <AudioControl bind:volume />
    </div>
</div>
