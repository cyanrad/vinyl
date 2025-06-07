<script lang="ts">
    import Vinyl from './Vinyl.svelte';
    import Controller from './Controller.svelte';
    import PlayerState from './PlayerState';
    import Scrubber from './Scrubber.svelte';

    // the overall state of the player coordinated with all components
    let playerState: PlayerState = $state(PlayerState.Paused);

    // aduio variables
    let audio: HTMLAudioElement = $state(new Audio('/music/avans-all_in.mp3'));
    let currentTime: number = $state(0);
    let duration: number = $state(0);

    // detect when a scrubber event happens
    let scrubberUpdated: boolean = $state(false);

    // handle audio playback based on player state
    $effect(() => {
        if (!audio) return;

        if (scrubberUpdated) {
            audio.currentTime = currentTime;
            scrubberUpdated = false; // reset the flag after updating
        }

        if (playerState === PlayerState.Playing) {
            audio.play();
        } else if (playerState === PlayerState.Paused) {
            audio.pause();
        }
        if (currentTime && duration && currentTime >= duration) {
            playerState = PlayerState.Paused; // stop playback if we reach the end
        }
    });

    // load duration when metadata is ready
    audio.addEventListener('loadedmetadata', () => {
        duration = audio.duration;
    });

    // update current time as it plays
    audio.addEventListener('timeupdate', () => {
        currentTime = audio.currentTime;
    });
</script>

<div class="relative">
    <img
        src="/base.svg"
        alt="Base"
        class="h-[620px] w-auto select-none pointer-events-none"
        draggable="false"
    />
    <Vinyl {playerState} coverUrl="../assets/phantasmagoria.jpg" />
    <div class="absolute top-[520px] left-[620px]">
        <Controller bind:playerState />
    </div>
    <div class="absolute top-[486px] left-[564px]">
        <Scrubber bind:currentTime bind:scrubberUpdated {duration} />
    </div>
</div>
