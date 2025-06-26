<script lang="ts">
    import PlayerState from "../PlayerState";
    import { getMovedAudioTime } from "../Audio";

    let {
        nextTrackSignal = $bindable(),
        playerState = $bindable(),
        currentTime = $bindable(),
        currTimeUpdated = $bindable(),
        duration,
        audio,
    } = $props();

    function playPause() {
        if (playerState === PlayerState.Paused && audio) {
            playerState = PlayerState.Playing;
        } else {
            playerState = PlayerState.Paused;
        }
    }

    // controller dimensions
    const controllerWidth = 230;
    const buttonSize = 30;
    const forwardButtonWidth = 47;
    const forwardButtonHeight = 30;

    // component images
    const previousTrackActiveImage = "controller/next-active.svg";
    const previousTrackInactiveImage = "controller/next-inactive.svg";
    const forwardTrackActiveImage = "controller/forward-active.svg";
    const forwardTrackInactiveImage = "controller/forward-inactive.svg";
    const playTrackInactiveImage = "controller/play-button-inactive.svg";
    const playTrackActiveImage = "controller/play-button-active.svg";
    const pauseTrackInactiveImage = "controller/pause-button-inactive.svg";
    const pauseTrackActiveImage = "controller/pause-button-active.svg";
    const nextTrackActiveImage = "controller/next-active.svg";
    const nextTrackInactiveImage = "controller/next-inactive.svg";

    // dynamically chaning pause/play based on player state
    let playButtonInactiveSrc = $derived.by(() => {
        if (playerState === PlayerState.Paused || !audio) {
            return playTrackInactiveImage;
        } else {
            return pauseTrackInactiveImage;
        }
    });
    let playButtonActiveSrc = $derived.by(() => {
        if (playerState === PlayerState.Paused || !audio) {
            return playTrackActiveImage;
        } else {
            return pauseTrackActiveImage;
        }
    });
</script>

<div class="flex flex-row items-center justify-between" style="width: {controllerWidth}px">
    <!-- preivous track -->
    <button
        class="relative inline-block cursor-pointer"
        style="width: {buttonSize}px"
        onclick={() => {
            nextTrackSignal = -1;
        }}
    >
        <img
            draggable="false"
            src={previousTrackInactiveImage}
            alt="play-button"
            class="w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 select-none"
            style="height: {buttonSize}px;"
        />
        <img
            draggable="false"
            src={previousTrackActiveImage}
            alt="play-button-active"
            class="absolute top-0 left-0 w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 select-none"
            style="height: {buttonSize}px;"
        />
    </button>

    <!-- rewind 5 s -->
    <button
        class="relative inline-block cursor-pointer mr-1"
        style="width: {forwardButtonWidth}px; height: {forwardButtonHeight}px"
        onclick={() => {
            if (!audio) return;

            currentTime = getMovedAudioTime(currentTime, duration, -5);
            currTimeUpdated = true;
        }}
    >
        <img
            draggable="false"
            src={forwardTrackInactiveImage}
            alt="play-button"
            class="w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 mx-auto select-none"
            style="height: {buttonSize}px;"
        />
        <img
            draggable="false"
            src={forwardTrackActiveImage}
            alt="play-button-active"
            class="absolute inset-0 m-auto w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 select-none"
            style="height: {buttonSize}px;"
        />
    </button>

    <!-- play/pause -->
    <button class="relative inline-block cursor-pointer" style="width: {buttonSize}px" onclick={playPause}>
        <img
            draggable="false"
            src={playButtonInactiveSrc}
            alt="play-button"
            class="w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 select-none"
            style="height: {buttonSize}px;"
        />
        <img
            draggable="false"
            src={playButtonActiveSrc}
            alt="play-button-active"
            class="absolute top-0 left-0 w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 select-none"
            style="height: {buttonSize}px;"
        />
    </button>

    <!-- forward 5 s -->
    <button
        class="relative inline-block cursor-pointer rotate-180 ml-1"
        style="width: {forwardButtonWidth}px; height: {forwardButtonHeight}px"
        onclick={() => {
            if (!audio) return;

            currentTime = getMovedAudioTime(currentTime, duration, 5);
            currTimeUpdated = true;
        }}
    >
        <img
            draggable="false"
            src={forwardTrackInactiveImage}
            alt="play-button"
            class="w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 mx-auto select-none"
            style="height: {buttonSize}px;"
        />
        <img
            draggable="false"
            src={forwardTrackActiveImage}
            alt="play-button-active"
            class="absolute inset-0 m-auto w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 select-none"
            style="height: {buttonSize}px;"
        />
    </button>

    <!-- next track -->
    <button
        class="relative inline-block cursor-pointer"
        style="width: {buttonSize}px"
        onclick={() => {
            nextTrackSignal = 1;
        }}
    >
        <img
            draggable="false"
            src={nextTrackInactiveImage}
            alt="play-button"
            class="w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 rotate-180 select-none"
            style="height: {buttonSize}px;"
        />
        <img
            draggable="false"
            src={nextTrackActiveImage}
            alt="play-button-active"
            class="absolute top-0 left-0 w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 rotate-180 select-none"
            style="height: {buttonSize}px;"
        />
    </button>
</div>
