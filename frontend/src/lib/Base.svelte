<script lang="ts">
    import Vinyl from './Vinyl.svelte';
    import Controller from './Controller.svelte';
    import PlayerState from './PlayerState';

    // the overall state of the player coordinated with all components
    let playerState: PlayerState = $state(PlayerState.Stopped);

    // aduio variables
    let audio: HTMLAudioElement = $state(new Audio('/music/avans-all_in.mp3'));
    let currentTime: number = $state(0);
    let duration: number = $state(0);

    // handle audio playback based on player state
    $effect(() => {
        console.log(`Current Time: ${currentTime}, Total Time: ${duration}`);
        if (playerState === PlayerState.Playing) {
            audio.play();
        } else if (playerState === PlayerState.Paused) {
            audio.pause();
        } else if (playerState === PlayerState.Stopped) {
            audio.pause();
            audio.currentTime = 0;
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
    <img src="/base.svg" alt="Base" class="h-[620px] w-auto" />
    <Vinyl {playerState} coverUrl="../assets/phantasmagoria.jpg" />
    <div class="absolute top-[520px] left-[620px]">
        <Controller bind:playerState />
    </div>
</div>
