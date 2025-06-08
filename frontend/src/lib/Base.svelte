<script lang="ts">
    import { onMount } from "svelte";

    import Vinyl from "./Vinyl.svelte";
    import Controller from "./Controller.svelte";
    import PlayerState from "./PlayerState";
    import Scrubber from "./Scrubber.svelte";

    import { getMovedAudioTime, getNewAudioTime } from "./audio";

    // the overall state of the player coordinated with all components
    let playerState: PlayerState = $state(PlayerState.Paused);

    // aduio variables
    let audio: HTMLAudioElement = $state(new Audio("/music/avans-all_in.mp3"));
    let currentTime: number = $state(0);
    let duration: number = $state(0);

    // detect when a scrubber event happens
    let currTimeUpdated: boolean = $state(false);

    // keyboard events
    onMount(() => {
        // play/pause events
        document.addEventListener("keydown", (event) => {
            if (event.code === "Space" || event.code === "KeyK") {
                event.preventDefault(); // Prevent page scroll
                playerState =
                    playerState === PlayerState.Playing
                        ? PlayerState.Paused
                        : PlayerState.Playing;
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
        }
    });

    // load duration when metadata is ready
    audio.addEventListener("loadedmetadata", () => {
        duration = audio.duration;
    });

    // update current time as it plays
    audio.addEventListener("timeupdate", () => {
        currentTime = audio.currentTime;
    });

    // base dimensions
    const baseHeight = 620;
    const controllerTop = 520;
    const controllerLeft = 620;
    const scrubberTop = 486;
    const scrubberLeft = 564;
</script>

<div class="relative">
    <img
        src="/base.svg"
        alt="Base"
        class="w-auto select-none pointer-events-none"
        style="height: {baseHeight}px;"
        draggable="false"
    />
    <Vinyl
        {playerState}
        coverUrl="../assets/phantasmagoria.jpg"
        {currentTime}
        {currTimeUpdated}
    />
    <div
        class="absolute"
        style="top: {controllerTop}px; left: {controllerLeft}px;"
    >
        <Controller
            bind:playerState
            bind:currentTime
            bind:currTimeUpdated
            {duration}
        />
    </div>
    <div class="absolute" style="top: {scrubberTop}px; left: {scrubberLeft}px;">
        <Scrubber bind:currentTime bind:currTimeUpdated {duration} />
    </div>
</div>
