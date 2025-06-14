<script lang="ts">
    import type { Track, Artist, Album } from "./api/Types";

    let { track } = $props();

    // diamentions
    const trackItemHeight = 60;
    const trackItemShadowLeft = 6;
    const trackItemShadowtop = 9;
</script>

<button class="relative group cursor-pointer">
    <!-- shadow -->
    <div
        class="absolute w-full bg-zinc-900 rounded-2xl group-hover:bg-emerald-800 transition-all duration-75 opacity-80"
        style="left: {trackItemShadowLeft}px; top: {trackItemShadowtop}px; height: {trackItemHeight}px"
    ></div>

    <!-- track item -->
    <div
        class="absolute flex top-0 left-0 bg-zinc-900 rounded-2xl w-full border-3 border-emerald-800 group-hover:border-emerald-400 transition-all duration-75"
        style="height: {trackItemHeight}px"
    >
        {track.title}
    </div>
</button>
