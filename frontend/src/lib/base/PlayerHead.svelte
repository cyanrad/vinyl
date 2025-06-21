<script lang="ts">
    import { onMount, untrack, getContext } from "svelte";

    import PlayerState from "../PlayerState";

    let { currentTime, currTimeUpdated, duration, playerState } = $props();
    let playerArmElement: HTMLImageElement | null = $state(null);

    // rotation values
    const minPausedArmRotation = 2;
    const minPlayingArmRotation = 7;
    const maxPlayingArmRotation = 27;
    const rotationPerSecond = $derived.by(() => {
        if (duration) {
            return (maxPlayingArmRotation - minPlayingArmRotation) / duration;
        }
        return 0;
    });
    let rotation: number = $state(minPausedArmRotation);

    // browser hidden state, requestAnimationFrame doesn't get computed when the browser is hidden
    let broswerHidden: boolean = $state(false);

    // animation frame ID for controlling the rotation on/off
    let animationId: number | null = null;

    // defaulting to 100 if not provided (much smoother than 60)
    // the animation logic relys on the monitor refresh rate, so we need to standarize the FPS
    const FPS: number = getContext("fps") || 100;
    const ROTATION_SPEED: number = 1.645; // not sure why the FPS is not enough but this is to make it spin once every minute

    // last timestamp to control the frame rate
    let lastTimestamp: DOMHighResTimeStamp | null = null;

    // Getting the vite dynamic image URLs from the assets folder
    onMount(async () => {
        if (playerState === PlayerState.Playing) {
            setRotationToTime(currentTime);
        } else if (playerState === PlayerState.Paused) {
            stopRotation();
        }
    });

    // effect to handle continuous rotation
    // not my favorite way to handle this, but couldn't find a better way
    // the animation is currently controlled by the `requestAnimationFrame` API
    // which is one of the few ways I can do the animation and retain the rotation variable.
    $effect(() => {
        if (playerState === PlayerState.Playing) {
            untrack(() => {
                setRotationToTime(currentTime);
                startRotation(null);
            });
        } else if (playerState === PlayerState.Paused || currentTime === 0) {
            untrack(() => {
                stopRotation();
            });
        }
    });

    // update the rotation if the current time changes due to skips or scrubbing
    $effect(() => {
        if (currTimeUpdated && playerState === PlayerState.Playing) {
            setRotationToTime(currentTime);
        }
    });

    // sets the rotation based on the current time of the audio
    function setRotationToTime(time: number) {
        if (playerState === PlayerState.Playing) {
            rotation = minPlayingArmRotation + time * rotationPerSecond;
        } else {
            rotation = minPausedArmRotation;
        }
        setRotation();
    }

    // will keep calling itself until we cancel the requestAnimationFrame
    // speed is controlled internall and by the FPS variable, this should prolly be changed later
    function startRotation(timestamp: DOMHighResTimeStamp | null) {
        animationId = requestAnimationFrame(startRotation);

        if (timestamp && lastTimestamp && timestamp - lastTimestamp < 1000 / FPS) {
            return; // Skip this frame to maintain FPS
        }

        rotation += (rotationPerSecond / FPS) * ROTATION_SPEED;
        setRotation();

        lastTimestamp = timestamp;
    }

    // cancels the rotation animation while retaining the current rotation state
    function stopRotation() {
        if (!animationId) return;

        cancelAnimationFrame(animationId);
        animationId = null;
        rotation = minPausedArmRotation;
        setRotation();
    }

    // sets the rotation of the vinyl and the cover image
    function setRotation() {
        if (!playerArmElement) return;
        playerArmElement.style.transform = `rotate(${rotation}deg)`;
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

    // player base dimensions
    const playerBaseHeight = 185;

    // player stand dimensions
    const playerStandHeight = 60;
    const playerStandTop = 58;
    const playerStandLeft = 50;

    // player arm dimensions & rotation origin
    const playerArmHeight = 430;
    const playerArmWidth = playerArmHeight * 0.553;
    const playerArmTop = 10;
    const playerArmLeft = -55;
    const playerArmOriginLeft = 67;
    const playerArmOriginTop = 20;

    // component images
    const playerHeadBaseImage = "player-head/base.svg";
    const playerHeadArmImage = "player-head/arm.svg";
    const playerHeadStandImage = "player-head/stand.svg";
</script>

<div class="relative">
    <!-- Player Head Base -->
    <img
        src={playerHeadBaseImage}
        alt=""
        class="w-auto select-none"
        style="height: {playerBaseHeight}px;"
        draggable="false"
    />

    <!-- Player Arm -->
    <!-- Stupid fucking thing refusing to work with height I have no clue why, had to surround it with a div-->
    <div
        class="absolute"
        style="height: {playerArmHeight}px; top: {playerArmTop}px; left: {playerArmLeft}px; width: {playerArmWidth}px;"
    >
        <img
            src={playerHeadArmImage}
            bind:this={playerArmElement}
            alt=""
            class="w-auto h-full absolute select-none"
            style="transform-origin: {playerArmOriginLeft}% {playerArmOriginTop}%; transform: rotate({rotation}deg);"
            draggable="false"
        />
    </div>

    <!-- Player Head Stand -->
    <img
        src={playerHeadStandImage}
        alt=""
        class="w-auto absolute select-none"
        style="height: {playerStandHeight}px; top: {playerStandTop}px; left: {playerStandLeft}px;"
        draggable="false"
    />
</div>
