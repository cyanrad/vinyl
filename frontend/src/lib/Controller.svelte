<script lang="ts">
    import PlayerState from './PlayerState';
    let { playerState = $bindable() } = $props();


    let playButtonInactiveSrc = $derived.by(() => {
        if (playerState === PlayerState.Paused) {
            return '/play-button-inactive.svg';
        } else {
            return '/pause-button-inactive.svg';
        }
    });
    let playButtonActiveSrc = $derived.by(() => {
        if (playerState === PlayerState.Paused) {
            return '/play-button-active.svg';
        } else {
            return '/pause-button-active.svg';
        }
    });

    function playPause() {
        if (playerState === PlayerState.Paused) {
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
</script>

<div class="flex flex-row items-center justify-between" style="width: {controllerWidth}px">
    <!-- preivous track -->
    <button class="relative inline-block cursor-pointer" style="width: {buttonSize}px">
        <img
            draggable="false"
            src="/next-inactive.svg"
            alt="play-button"
            class="w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 select-none"
            style="height: {buttonSize}px;"
        />
        <img
            draggable="false"
            src="/next-active.svg"
            alt="play-button-active"
            class="absolute top-0 left-0 w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 select-none"
            style="height: {buttonSize}px;"
        />
    </button>

    <!-- rewind 5 s -->
    <button class="relative inline-block cursor-pointer mr-1" style="width: {forwardButtonWidth}px; height: {forwardButtonHeight}px">
        <img
            draggable="false"
            src="/forward-inactive.svg"
            alt="play-button"
            class="w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 mx-auto select-none"
            style="height: {buttonSize}px;"
        />
        <img
            draggable="false"
            src="/forward-active.svg"
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
    <button class="relative inline-block cursor-pointer rotate-180 ml-1" style="width: {forwardButtonWidth}px; height: {forwardButtonHeight}px">
        <img
            draggable="false"
            src="/forward-inactive.svg"
            alt="play-button"
            class="w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 mx-auto select-none"
            style="height: {buttonSize}px;"
        />
        <img
            draggable="false"
            src="/forward-active.svg"
            alt="play-button-active"
            class="absolute inset-0 m-auto w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 select-none"
            style="height: {buttonSize}px;"
        />
    </button>

    <!-- next track -->
    <button class="relative inline-block cursor-pointer" style="width: {buttonSize}px">
        <img
            draggable="false"
            src="/next-inactive.svg"
            alt="play-button"
            class="w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 rotate-180 select-none"
            style="height: {buttonSize}px;"
        />
        <img
            draggable="false"
            src="/next-active.svg"
            alt="play-button-active"
            class="absolute top-0 left-0 w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 rotate-180 select-none"
            style="height: {buttonSize}px;"
        />
    </button>
</div>
