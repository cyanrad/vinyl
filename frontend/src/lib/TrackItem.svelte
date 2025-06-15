<script lang="ts">
    import { setContext } from "svelte";
    import { fly } from "svelte/transition";

    import { generateTrackCoverUrl } from "./api/Tracks";

    let { track, index, activeTrackIndex = $bindable() } = $props();

    let borderColor = $derived(activeTrackIndex === index ? "border-emerald-400" : "border-zinc-700");
    let textColor = $derived(activeTrackIndex === index ? "text-emerald-400" : "text-emerald-600");

    function selectTrack() {
        activeTrackIndex = index;
    }

    // diamentions
    const trackItemHeight = 60;
    const trackItemShadowLeft = 6;
    const trackItemShadowtop = 9;
</script>

<button
    class="relative group cursor-pointer hover:scale-103 transition-all duration-50"
    in:fly|global={{ duration: 500, delay: index * 100, y: 100 }}
    onclick={selectTrack}
>
    <!-- shadow -->
    <div
        class="absolute w-full bg-zinc-900 rounded-2xl transition-all duration-150 opacity-80"
        style="left: {trackItemShadowLeft}px; top: {trackItemShadowtop}px; height: {trackItemHeight}px"
    ></div>

    <!-- track item -->
    <div
        class="absolute flex top-0 left-0 bg-zinc-900 rounded-2xl w-full items-center border-4 p-1 {borderColor} group-hover:border-emerald-400 transition-all duration-75"
        style="height: {trackItemHeight}px"
    >
        <img src={generateTrackCoverUrl(track)} alt="" class="h-full w-auto rounded-xl aspect-square opacity-80 mr-2" />
        <div class="flex flex-col">
            <span class="font-bold font-sm quicksand {textColor} group-hover:text-emerald-400" style="font-size: 13px"
                >{track.title}</span
            >
            <!-- <span class="font-bold quicksand text-zinc-400" style="font-size: 12px">{}</span> -->
        </div>
    </div>
</button>
