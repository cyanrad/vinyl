<script lang="ts">
    import { onMount, untrack, getContext } from 'svelte';
    import PlayerState from './PlayerState';

    // vinyl record dimensions and positioning
    const vinylRecordTop = 65;
    const vinylRecordLeft = 50;
    const vinylRecordHeight = 440;

    // album cover dimensions and positioning
    const albumCoverTop = 200;
    const albumCoverLeft = 185;
    const albumCoverHeight = 170;

    // vinyl center piece dimensions and positioning
    const vinylCenterTop = 268;
    const vinylCenterLeft = 254;
    const vinylCenterHeight = 32;

    // the rotation of the vinyl image
    let rotation: number = $state(0);

    // variables to hold to the image elements for rotation logic
    let vinylElement: HTMLImageElement | null = $state(null); // Reference to the image element for rotation
    let coverElement: HTMLImageElement | null = $state(null); // Reference to the cover image element

    // animation frame ID for controlling the rotation on/off
    let animationId: number | null = null;

    // playerState of the vinyl player, can be 'spinning' or 'paused'
    // coverUrl is the URL of the cover image of the track
    let { playerState, coverUrl } = $props();

    // defaulting to 100 if not provided (much smoother than 60)
    // the animation logic relys on the monitor refresh rate, so we need to standarize the FPS
    const FPS: number = getContext('fps') || 100;
    // last timestamp to control the frame rate
    let lastTimestamp: DOMHighResTimeStamp | null = null;

    // Getting the vite dynamic image URLs from the assets folder
    onMount(async () => {
        if (coverUrl !== undefined) {
            const module = await import(coverUrl);
            coverUrl = module.default;
        }
    });

    // effect to handle continuous rotation
    // not my favorite way to handle this, but couldn't find a better way
    // the animation is currently controlled by the `requestAnimationFrame` API
    // which is one of the few ways I can do the animation and retain the rotation variable.
    $effect(() => {
        if (playerState === PlayerState.Playing) {
            untrack(() => {
                startRotation(null);
            });
        } else if (playerState === PlayerState.Paused) {
            untrack(() => {
                stopRotation();
            });
        }
    });

    // will keep calling itself until we cancel the requestAnimationFrame
    // speed is controlled internall and by the FPS variable, this should prolly be changed later
    // todo - the rotation of the disk is not computed when the browser is not open
    //        considering that we have a framerate it can be possible to compute the rotation
    function startRotation(timestamp: DOMHighResTimeStamp | null) {
        animationId = requestAnimationFrame(startRotation);

        if (timestamp && lastTimestamp && timestamp - lastTimestamp < 1000 / FPS) {
            return; // Skip this frame to maintain FPS
        }

        rotation += 0.1;
        if (vinylElement) {
            vinylElement.style.transform = `rotate(${rotation}deg)`;
        }

        if (coverElement) {
            coverElement.style.transform = `rotate(${rotation}deg)`;
        }

        lastTimestamp = timestamp;
    }

    // cancels the rotation animation while retaining the current rotation state
    function stopRotation() {
        if (animationId) {
            cancelAnimationFrame(animationId);
            animationId = null;
        }
    }

    // resets the rotation to 0 degrees
    function resetRotation() {
        if (vinylElement) {
            vinylElement.style.transform = `rotate(0deg)`;
        }

        if (coverElement) {
            coverElement.style.transform = `rotate(0deg)`;
        }
        rotation = 0;
    }
</script>

<!-- Should've probably surrounded them with a div or something -->
<img
    draggable="false"
    bind:this={vinylElement}
    src="/vinyl.png"
    alt="Vinyl"
    class="absolute w-auto select-none pointer-events-none"
    style="top: {vinylRecordTop}px; left: {vinylRecordLeft}px; height: {vinylRecordHeight}px;"
/>
{#if coverUrl}
    <img
        draggable="false"
        bind:this={coverElement}
        src={coverUrl}
        alt=""
        class="absolute w-auto z-10 rounded-full object-cover select-none pointer-events-none"
        style="top: {albumCoverTop}px; left: {albumCoverLeft}px; height: {albumCoverHeight}px;"
    />
{/if}
<img
    draggable="false"
    src="/vinyl-center.png"
    alt="Vinyl Center"
    class="absolute w-auto z-20 select-none pointer-events-none"
    style="top: {vinylCenterTop}px; left: {vinylCenterLeft}px; height: {vinylCenterHeight}px;"
/>

<!-- if removed will cause imports to fail due to seemingly how the animation is done -->
<style>
</style>
