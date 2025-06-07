<script lang="ts">
    import { onMount, untrack, getContext } from 'svelte';
    import PlayerState from './PlayerState';

    // the rotation of the vinyl image
    let rotation: number = $state(0);

    // variables to hold to the image elements for rotation logic
    let vinylElement: HTMLImageElement | null = $state(null); // Reference to the image element for rotation
    let coverElement: HTMLImageElement | null = $state(null); // Reference to the cover image element

    // animation frame ID for controlling the rotation on/off
    let animationId: number | null = null;

    // playerState of the vinyl player, can be 'spinning', 'stopped', or 'paused'
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
        } else if (playerState === PlayerState.Stopped) {
            untrack(() => {
                stopRotation();
                resetRotation();
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
    bind:this={vinylElement}
    src="/vinyl.png"
    alt="Vinyl"
    class="absolute top-[65px] left-[50px] h-[440px] w-auto"
/>
{#if coverUrl}
    <img
        bind:this={coverElement}
        src={coverUrl}
        alt="Cover"
        class="absolute top-[200px] left-[185px] h-[170px] w-auto z-10 rounded-full object-cover"
    />
{/if}
<img
    src="/vinyl-center.png"
    alt="Vinyl Center"
    class="absolute top-[268px] left-[254px] h-[32px] w-auto z-20"
/>

<!-- Control buttons to be deleted -->
<!-- <div class="controls"> -->
<!--     <button -->
<!--         onmousedown={() => { -->
<!--             playerState = 'spinning'; -->
<!--         }}>Continuous Spin</button -->
<!--     > -->
<!--     <button -->
<!--         onmousedown={() => { -->
<!--             playerState = 'stopped'; -->
<!--         }}>stop Spin</button -->
<!--     > -->
<!--     <button -->
<!--         onmousedown={() => { -->
<!--             playerState = 'paused'; -->
<!--         }}>puase Spin</button -->
<!--     > -->
<!-- </div> -->

<!-- if removed will cause imports to fail due to seemingly how the animation is done -->
<style>
</style>
