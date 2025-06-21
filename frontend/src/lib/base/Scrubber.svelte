<script lang="ts">
    // audio variables
    let { currentTime = $bindable(), currTimeUpdated = $bindable(), duration } = $props();

    // scrub & dragging logic
    let isDragging = $state(false);
    let dragProgress = $state(0);
    let scrubberElement: HTMLButtonElement | null = $state(null);
    let scrubProgress = $derived.by(() => {
        if (isDragging) {
            return dragProgress;
        } else if (duration && currentTime) {
            // condition prevents full highlight if audio is not loaded
            return currentTime / duration;
        } else {
            return 0;
        }
    });

    let timeDisplayCurrentTime: number = $derived.by(() => {
        if (isDragging) {
            return dragProgress * duration;
        } else if (duration && currentTime) {
            // condition prevents full highlight if audio is not loaded
            return currentTime;
        } else {
            return 0;
        }
    });

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
        // Calculate final time and call onSeek
        currentTime = dragProgress * duration;
        currTimeUpdated = true;

        isDragging = false;

        // Remove global listeners
        document.removeEventListener("mousemove", handleMouseMove);
        document.removeEventListener("mouseup", handleMouseUp);
    }

    function updateDragProgress(event: MouseEvent) {
        if (!scrubberElement) return;

        const rect = scrubberElement.getBoundingClientRect();
        const x = event.clientX - rect.left;

        // if the click is outside the scrubber bounds, clamp it & get the for inbetween progress
        const progressPerc = x / rect.width;
        dragProgress = Math.max(0, Math.min(1, progressPerc));
    }

    function convertToMMSS(time: number | null) {
        if (!time) return "0:00";

        const minutes = Math.floor(time / 60);
        const seconds = Math.floor(time % 60);
        return `${minutes}:${seconds.toString().padStart(2, "0")}`;
    }

    // size controls
    const scrubberWidth = 330;
    const scrubberHeight = 6;
    const playheadSize = 16;
    const playheadTopOffset = -5;
    const playheadInnerSize = 12;
    const playheadInnerOffset = 2;
</script>

<div class="flex flex-col justify-center" style="width: {scrubberWidth}px">
    <!-- scrub bar -->
    <button
        bind:this={scrubberElement}
        class="w-full rounded-full bg-zinc-900 relative cursor-pointer"
        style="height: {scrubberHeight}px"
        onmousedown={handleMouseDown}
    >
        <!-- scrub bar highlight -->
        <div
            class="top-0 left-0 bg-emerald-400 rounded-full pointer-events-none"
            style="width: {scrubProgress * 100 + 2}%; height: {scrubberHeight}px"
        ></div>

        <!-- playhead -->
        <div
            class="bg-emerald-400 rounded-full absolute cursor-grab transition-transform hover:scale-110"
            class:cursor-grabbing={isDragging}
            style="left: {scrubProgress *
                100}%; width: {playheadSize}px; height: {playheadSize}px; top: {playheadTopOffset}px; transform: translateX(-50%);"
            onmousedown={handleMouseDown}
            role="slider"
            tabindex="0"
            aria-valuemin="0"
            aria-valuemax={duration}
            aria-valuenow={currentTime}
        >
            <div
                class="bg-zinc-900 rounded-full absolute pointer-events-none"
                style="width: {playheadInnerSize}px; height: {playheadInnerSize}px; top: {playheadInnerOffset}px; left: {playheadInnerOffset}px;"
            ></div>
        </div>
    </button>

    <!-- time display -->
    <div class="flex justify-between text-xs text-emerald-400 text-center mt-2 font-quicksand font-bold select-none">
        <span>{convertToMMSS(timeDisplayCurrentTime)}</span>
        <span>{convertToMMSS(duration)}</span>
    </div>
</div>
