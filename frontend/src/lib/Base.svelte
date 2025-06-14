<script lang="ts">
    import { onMount } from "svelte";

    // components
    import Vinyl from "./Vinyl.svelte";
    import Controller from "./Controller.svelte";
    import Scrubber from "./Scrubber.svelte";
    import PlayerHead from "./PlayerHead.svelte";
    import AudioControl from "./AudioControl.svelte";
    import DisplayMonitor from "./DisplayMonitor.svelte";

    // state
    import PlayerState from "./PlayerState";
    import { getMovedAudioTime, getNewAudioTime } from "./Audio";

    // API
    import { generateTrackAudioUrl } from "./api/Tracks";

    let { activeTrack, activeArtist, activeAlbum } = $props();

    // get the track audio url
    let trackAudio: string | null = $derived(activeTrack ? generateTrackAudioUrl(activeTrack) : null);

    // the overall state of the player coordinated with all components
    let playerState: PlayerState = $state(PlayerState.Paused);

    // audio variables
    let audio: HTMLAudioElement | null = $derived(trackAudio ? new Audio(trackAudio) : null);
    let currentTime: number = $state(0);
    let duration: number = $state(0);
    let volume: number = $state(0.5);

    // flag for when a scrubber event happens
    let currTimeUpdated: boolean = $state(false);

    // keyboard events
    onMount(async () => {
        // play/pause events
        document.addEventListener("keydown", (event) => {
            if (!audio) return;

            if (event.code === "Space" || event.code === "KeyK") {
                event.preventDefault(); // Prevent page scroll
                playerState = playerState === PlayerState.Playing ? PlayerState.Paused : PlayerState.Playing;
            }
        });

        // forward/rewind events
        document.addEventListener("keydown", (event) => {
            if (!audio) return;

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

    // display screen dimensions
    const displayScreenTop = 276;
    const displayScreenLeft = 562;

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

    <!-- display screen -->
    <div class="absolute z-40" style="top: {displayScreenTop}px; left: {displayScreenLeft}px; ">
        <DisplayMonitor {playerState} {activeTrack} {activeArtist} {activeAlbum} {audio} />
    </div>

    <!-- player head -->
    <div class="absolute z-30" style="top: {playerHeadTop}px; left: {playerHeadLeft}px;">
        <PlayerHead {currentTime} {duration} {playerState} {currTimeUpdated} />
    </div>

    <!-- vinyl -->
    <Vinyl {playerState} {activeTrack} {activeAlbum} {activeArtist} {currentTime} {currTimeUpdated} />

    <!-- controller -->
    <div class="absolute" style="top: {controllerTop}px; left: {controllerLeft}px;">
        <Controller bind:playerState bind:currentTime bind:currTimeUpdated {duration} {audio} />
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
