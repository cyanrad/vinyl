<script lang="ts">
    import { fly } from "svelte/transition";

    import TrackItem from "./TrackItem.svelte";

    let isSideBarOpen = $state(false);

    function toggleSideBar() {
        isSideBarOpen = !isSideBarOpen;
    }

    const playlistBarWidth = 60;
    const infoBarWidth = 310;
    const innerInfoBarWidth = 290;
    const innerInfoBarShadowWidth = 8;
</script>

<aside class="fixed left-0 top-0 h-screen z-10">
    <!-- wrapping stuff in a div because for some reason it fails if there is more than one component in the aside element -->
    <div class="flex items-start h-full w-full">
        <!-- playlist bar -->
        <div class="flex flex-col items-start h-full bg-zinc-900 p-2 z-20" style="width: {playlistBarWidth}px;">
            <button class="cursor-pointer w-full h-auto" onclick={toggleSideBar}>
                <img src="/vinyl-icon.svg" alt="" class="w-full h-full" />
            </button>
        </div>

        <!-- info bar -->
        {#if isSideBarOpen}
            <div class="relative h-full z-10" transition:fly={{ x: -infoBarWidth, duration: 300, opacity: 90 }}>
                <!-- outer -->
                <div
                    class="absolute top-0 left-0 h-full bg-zinc-700 rounded-r-xl"
                    style="width: {infoBarWidth}px"
                ></div>

                <!-- inner bar shadow -->
                <div
                    class="absolute top-[2%] h-[97%] bg-zinc-900 rounded-r-xl"
                    style="width: {innerInfoBarWidth + innerInfoBarShadowWidth}px;"
                ></div>

                <!-- inner bar -->
                <div
                    class="absolute top-[0.8%] h-[97%] bg-zinc-800 rounded-r-xl pt-2"
                    style="width: {innerInfoBarWidth}px"
                >
                    <div
                        class="flex flex-col gap-20 h-[calc(100%-10px)] overflow-y-auto p-4 rounded-4xl scroll-smooth scrollbar-none"
                    >
                        <TrackItem />
                        <TrackItem />
                        <TrackItem />
                        <TrackItem />
                        <TrackItem />
                        <TrackItem />
                        <TrackItem />
                        <TrackItem />
                        <TrackItem />
                        <TrackItem />
                        <TrackItem />
                    </div>
                </div>
            </div>
        {/if}
    </div>
</aside>
