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
</script>

<div class="w-[230px] flex flex-row items-center justify-between">
    <!-- preivous track -->
    <button class="relative inline-block cursor-pointer w-[30px]">
        <img
            draggable="false"
            src="/next-inactive.svg"
            alt="play-button"
            class="h-[30px] w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 select-none"
        />
        <img
            draggable="false"
            src="/next-active.svg"
            alt="play-button-active"
            class="absolute top-0 left-0 h-[30px] w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 select-none"
        />
    </button>

    <!-- rewind 5 s -->
    <button class="relative inline-block cursor-pointer w-[47px] h-[30px] mr-1">
        <img
            draggable="false"
            src="/forward-inactive.svg"
            alt="play-button"
            class="h-[30px] w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 mx-auto select-none"
        />
        <img
            draggable="false"
            src="/forward-active.svg"
            alt="play-button-active"
            class="absolute inset-0 m-auto h-[30px] w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 select-none"
        />
    </button>

    <!-- play/pause -->
    <button class="relative inline-block cursor-pointer w-[30px]" onclick={playPause}>
        <img
            draggable="false"
            src={playButtonInactiveSrc}
            alt="play-button"
            class="h-[30px] w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 select-none"
        />
        <img
            draggable="false"
            src={playButtonActiveSrc}
            alt="play-button-active"
            class="absolute top-0 left-0 h-[30px] w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 select-none"
        />
    </button>

    <!-- forward 5 s -->
    <button class="relative inline-block cursor-pointer w-[47px] h-[30px] rotate-180 ml-1">
        <img
            draggable="false"
            src="/forward-inactive.svg"
            alt="play-button"
            class="h-[30px] w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 mx-auto select-none"
        />
        <img
            draggable="false"
            src="/forward-active.svg"
            alt="play-button-active"
            class="absolute inset-0 m-auto h-[30px] w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 select-none"
        />
    </button>

    <!-- next track -->
    <button class="relative inline-block cursor-pointer w-[30px]">
        <img
            draggable="false"
            src="/next-inactive.svg"
            alt="play-button"
            class="h-[30px] w-auto transition-opacity duration-300 ease-in-out hover:opacity-0 rotate-180 select-none"
        />
        <img
            draggable="false"
            src="/next-active.svg"
            alt="play-button-active"
            class="absolute top-0 left-0 h-[30px] w-auto opacity-0 transition-opacity duration-300 ease-in-out hover:opacity-100 rotate-180 select-none"
        />
    </button>
</div>
