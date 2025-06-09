<script lang="ts">
    import { onMount } from "svelte";

    // audio variables
    let { volume = $bindable() } = $props();

    let audioScrubberElement: HTMLButtonElement | null = $state(null);

    // scrub & dragging logic
    let isDragging = $state(false);

    // setting up event listeners and flags
    function handleMouseDown(event: MouseEvent) {
        isDragging = true;
        updateDragProgress(event);

        // add global mouse listeners
        document.addEventListener("mousemove", handleMouseMove);
        document.addEventListener("mouseup", handleMouseUp);

        // prevent text selection
        event.preventDefault();
    }

    function handleMouseMove(event: MouseEvent) {
        updateDragProgress(event);
    }

    function handleMouseUp(_: MouseEvent) {
        isDragging = false;

        // Remove global listeners
        document.removeEventListener("mousemove", handleMouseMove);
        document.removeEventListener("mouseup", handleMouseUp);
    }

    function updateDragProgress(event: MouseEvent) {
        if (!audioScrubberElement) return;

        const rect = audioScrubberElement.getBoundingClientRect();
        const y = rect.bottom - event.clientY;

        // if the click is outside the scrubber bounds, clamp it & get the for inbetween progress
        const progressPerc = y / rect.height;
        volume = Math.max(0, Math.min(1, progressPerc));
    }

    const AudioScrubberWidth = 6;
    const AudioScrubberHeight = 164;

    // playhead
    const AudioPlayheadWidth = 26;
    const AudioPlayheadHeight = 8;
    const AudioPlayheadLeftOffset = (-AudioPlayheadWidth + AudioScrubberWidth) / 2;
    const AudioPlayheadInnerWidth = 20;
    const AudioPlayheadInnerHeight = 4;
</script>

<div class="flex flex-col justify-center">
    <!-- scrub bar -->
    <button
        bind:this={audioScrubberElement}
        class="rounded-full bg-zinc-900 relative cursor-pointer"
        style="width: {AudioScrubberWidth}px; height: {AudioScrubberHeight}px;"
        onmousedown={handleMouseDown}
    >
        <!-- scrub bar highlight -->
        <div
            class="bg-emerald-400 rounded-full pointer-events-none absolute"
            style="width: {AudioScrubberWidth}px; height: {volume * 100 + 2}%; bottom: 0%;"
        ></div>

        <!-- playhead -->
        <div
            class="bg-emerald-400 rounded-full absolute cursor-grab transition-transform hover:scale-110 flex justify-center items-center"
            class:cursor-grabbing={isDragging}
            style="bottom: {volume *
                100}%; width: {AudioPlayheadWidth}px; height: {AudioPlayheadHeight}px; left: {AudioPlayheadLeftOffset}px;"
            onmousedown={handleMouseDown}
            role="slider"
            tabindex="0"
            aria-valuemin="0"
            aria-valuemax="1"
            aria-valuenow={volume}
        >
            <div
                class="bg-zinc-900 rounded-full pointer-events-none"
                style="width: {AudioPlayheadInnerWidth}px; height: {AudioPlayheadInnerHeight}px;"
            ></div>
        </div>
    </button>
</div>
