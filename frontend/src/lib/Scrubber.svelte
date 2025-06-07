<script lang="ts">
    // audio variables
    let { currentTime = $bindable(), scrubberUpdated = $bindable(), duration } = $props();

    // size controls
    const scrubberWidth = 330;
    const scrubberHeight = 6;
    const playheadSize = 16;
    const playheadTopOffset = -5;
    const playheadInnerSize = 12;
    const playheadInnerOffset = 2;

    // scrub & dragging logic logic
    let isDragging = $state(false);
    let dragProgress = $state(0);
    let scrubberElement: HTMLButtonElement | null = $state(null);
    let scrubProgress = $derived(isDragging ? dragProgress : (currentTime / duration) * 100);

    // setting up event listeners and flags
    function handleMouseDown(event: MouseEvent) {
        isDragging = true;
        updateDragProgress(event);

        // add global mouse listeners
        document.addEventListener('mousemove', handleMouseMove);
        document.addEventListener('mouseup', handleMouseUp);

        // prevent text selection
        event.preventDefault();
    }

    function handleMouseMove(event: MouseEvent) {
        updateDragProgress(event);
    }

    function handleScrubberClick(event: MouseEvent) {
        if (isDragging) return;

        updateDragProgress(event);
        currentTime = (dragProgress / 100) * duration;
        scrubberUpdated = true;
    }

    function handleMouseUp(_: MouseEvent) {
        // Calculate final time and call onSeek
        currentTime = (dragProgress / 100) * duration;
        scrubberUpdated = true;

        isDragging = false;

        // Remove global listeners
        document.removeEventListener('mousemove', handleMouseMove);
        document.removeEventListener('mouseup', handleMouseUp);
    }

    function updateDragProgress(event: MouseEvent) {
        if (!scrubberElement) return;

        const rect = scrubberElement.getBoundingClientRect();
        const x = event.clientX - rect.left;

        // if the click is outside the scrubber bounds, clamp it & get the for inbetween progress
        const progressPerc = (x / rect.width) * 100;
        dragProgress = Math.max(0, Math.min(100, progressPerc));
    }
</script>

<div class="flex flex-col justify-center" style="width: {scrubberWidth}px">
    <!-- scrub bar -->
    <button
        bind:this={scrubberElement}
        class="w-full rounded-full bg-zinc-900 relative cursor-pointer"
        style="height: {scrubberHeight}px"
        onclick={handleScrubberClick}
    >
        <!-- scrub bar highlight -->
        <div
            class="top-0 left-0 bg-emerald-400 rounded-full pointer-events-none"
            style="width: {scrubProgress + 2}%; height: {scrubberHeight}px"
        ></div>

        <!-- playhead -->
        <div
            class="bg-emerald-400 rounded-full absolute cursor-grab transition-transform hover:scale-110"
            class:cursor-grabbing={isDragging}
            style="left: {scrubProgress}%; width: {playheadSize}px; height: {playheadSize}px; top: {playheadTopOffset}px; transform: translateX(-50%);"
            onmousedown={handleMouseDown}
            role="slider"
            tabindex="0"
            aria-valuemin="0"
            aria-valuemax={duration}
            aria-valuenow={isDragging ? (dragProgress / 100) * duration : currentTime}
        >
            <div
                class="bg-zinc-900 rounded-full absolute pointer-events-none"
                style="width: {playheadInnerSize}px; height: {playheadInnerSize}px; top: {playheadInnerOffset}px; left: {playheadInnerOffset}px;"
            ></div>
        </div>
    </button>
</div>
