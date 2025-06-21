<script lang="ts">
    import { onMount, untrack, getContext } from "svelte";
    import PlayerState from "../PlayerState";

    // API
    import {
        generateTrackItemCoverUrl,
        generateTrackItemAlbumCoverUrl,
        generateTrackItemArtistImageUrl,
    } from "../api/TrackItems";

    // We need information about the general play state, audio timeing and track data
    let { playerState, activeTrack, currentTime, currTimeUpdated } = $props();

    // the rotation of the vinyl image
    let rotation: number = $state(0);

    // browser hidden state, requestAnimationFrame doesn't get computed when the browser is hidden
    let broswerHidden: boolean = $state(false);

    // variables to hold to the image elements for rotation logic
    let vinylElement: HTMLImageElement | null = $state(null);
    let coverElement: HTMLImageElement | null = $state(null);

    // image of the vinyl cover
    let coverImage: string | null = $derived.by(() => {
        if (activeTrack && activeTrack.cover) {
            return generateTrackItemCoverUrl(activeTrack);
        } else if (activeTrack && activeTrack.albumCover) {
            return generateTrackItemAlbumCoverUrl(activeTrack);
        } else if (activeTrack && activeTrack.artistImages[0]) {
            return generateTrackItemArtistImageUrl(activeTrack, 0);
        } else {
            return null;
        }
    });

    // animation frame ID for controlling the rotation on/off
    let animationId: number | null = null;

    // defaulting to 100 if not provided (much smoother than 60)
    // the animation logic relys on the monitor refresh rate, so we need to standarize the FPS
    const FPS: number = getContext("fps") || 100;
    const ROTATION_SPEED: number = 1.645; // not sure why the FPS is not enough but this is to make it spin once every minute
    // last timestamp to control the frame rate
    let lastTimestamp: DOMHighResTimeStamp | null = null;

    // getting the vite dynamic image URLs from the assets folder
    onMount(async () => {
        setRotationToTime(currentTime);
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

    // update the rotation if the current time changes due to skips or scrubbing
    $effect(() => {
        if (currTimeUpdated) {
            setRotationToTime(currentTime);
        }
    });

    // sets the rotation of the vinyl and the cover image based on the current time of the audio
    function setRotationToTime(time: number) {
        rotation = time * (360 / 60);
        setRotation(rotation);
    }

    // will keep calling itself until we cancel the requestAnimationFrame
    // speed is controlled internall and by the FPS variable, this should prolly be changed later
    function startRotation(timestamp: DOMHighResTimeStamp | null) {
        animationId = requestAnimationFrame(startRotation);

        if (timestamp && lastTimestamp && timestamp - lastTimestamp < 1000 / FPS) {
            return; // Skip this frame to maintain FPS
        }

        rotation += (360 / 60 / FPS) * ROTATION_SPEED; // one rotation per minute
        setRotation(rotation);

        lastTimestamp = timestamp;
    }

    // cancels the rotation animation while retaining the current rotation state
    function stopRotation() {
        if (animationId) {
            cancelAnimationFrame(animationId);
            animationId = null;
        }
    }

    // sets the rotation of the vinyl and the cover image
    function setRotation(rotation: number) {
        if (vinylElement) {
            vinylElement.style.transform = `rotate(${rotation}deg)`;
        }

        if (coverElement) {
            coverElement.style.transform = `rotate(${rotation}deg)`;
        }
    }

    // when the browser is hidden, the animation is paused and rotation does not change
    // when the browser is visible again we want to resume the rotation based on current time
    document.addEventListener("visibilitychange", () => {
        if (document.visibilityState === "hidden") {
            broswerHidden = true;
        } else if (document.visibilityState === "visible" && broswerHidden) {
            setRotationToTime(currentTime);
            broswerHidden = false;
        }
    });

    // vinyl record dimensions and positioning
    const vinylRecordTop = 65;
    const vinylRecordLeft = 50;
    const vinylRecordHeight = 440;

    // album cover dimensions and positioning
    const albumCoverTop = 200;
    const albumCoverLeft = 185;
    const albumCoverLength = 170;

    // vinyl center piece dimensions and positioning
    const vinylCenterTop = 268;
    const vinylCenterLeft = 254;
    const vinylCenterHeight = 32;

    // component images
    const vinylDiskImage = "vinyl/disk.svg";
    const vinylCenterImage = "vinyl/center.svg";
</script>

<!-- Should've probably surrounded them with a div or something -->
<img
    draggable="false"
    bind:this={vinylElement}
    src={vinylDiskImage}
    alt="Vinyl"
    class="absolute w-auto select-none pointer-events-none"
    style="top: {vinylRecordTop}px; left: {vinylRecordLeft}px; height: {vinylRecordHeight}px;"
/>
{#if coverImage}
    <img
        draggable="false"
        bind:this={coverElement}
        src={coverImage}
        alt=""
        class="absolute w-auto z-10 rounded-full object-cover select-none pointer-events-none"
        style="top: {albumCoverTop}px; left: {albumCoverLeft}px; height: {albumCoverLength}px; width: {albumCoverLength}px;"
    />
    <div
        class="absolute w-auto z-20 rounded-full object-cover select-none pointer-events-none opacity-15 bg-zinc-900"
        style="top: {albumCoverTop - 1}px; left: {albumCoverLeft - 1}px; height: {albumCoverLength +
            1}px; width: {albumCoverLength + 1}px;"
    ></div>
{/if}
<img
    draggable="false"
    src={vinylCenterImage}
    alt="Vinyl Center"
    class="absolute w-auto z-10 select-none pointer-events-none"
    style="top: {vinylCenterTop}px; left: {vinylCenterLeft}px; height: {vinylCenterHeight}px;"
/>

<!-- if removed will cause imports to fail due to seemingly how the animation is done -->
<style>
</style>
