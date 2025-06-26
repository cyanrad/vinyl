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
    import PlayerState from "../PlayerState";
    import { getMovedAudioTime, getNewAudioTime } from "../Audio";

    let {
        activeTrack,
        nextTrackSignal = $bindable(),
        duration,
        playerState = $bindable(),
        audio = $bindable(),
        currentTime = $bindable(),
        currTimeUpdated = $bindable(),
    } = $props();

    // audio variables
    let volume: number = $state(0.05);

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
        <DisplayMonitor {playerState} {activeTrack} {audio} />
    </div>

    <!-- player head -->
    <div class="absolute z-30" style="top: {playerHeadTop}px; left: {playerHeadLeft}px;">
        <PlayerHead {currentTime} {duration} {playerState} {currTimeUpdated} />
    </div>

    <!-- vinyl -->
    <Vinyl {playerState} {activeTrack} {currentTime} {currTimeUpdated} />

    <!-- controller -->
    <div class="absolute" style="top: {controllerTop}px; left: {controllerLeft}px;">
        <Controller bind:playerState bind:currentTime bind:currTimeUpdated bind:nextTrackSignal {duration} {audio} />
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
