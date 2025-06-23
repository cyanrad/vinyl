<script lang="ts">
    import { onMount } from "svelte";

    import { generateTrackItemAlbumCoverUrl, generateTrackItemArtistImageUrl } from "../api/TrackItems";
    import PlayerState from "../PlayerState";

    let { playerState, activeTrack, audio } = $props();

    // text
    let trackTitle: string = $derived(activeTrack?.title || "");
    let artistName: string = $derived(activeTrack?.artistNames[0] || "");
    let albumTitle: string = $derived(activeTrack?.albumTitle || "404\nNo Album\n:D");

    let playerStateText: string = $derived.by(() => {
        if (!audio) return "NO TRACK SELECTED :[";
        return playerState === PlayerState.Playing ? "CURRENTLY PLAYING :)" : "PAUSED :O";
    });

    // images
    let albumCover: string | null = $derived(generateTrackItemAlbumCoverUrl(activeTrack));
    let artistImage: string | null = $derived(generateTrackItemArtistImageUrl(activeTrack));

    // flicker animation
    const flickerFPS = 24;
    const flickerOpacityMin = 0.18;
    const flickerOpacityMax = 0.3;
    const flickerSkew = 1;
    const textFlickerModifier = 0.7; // this gets added on top of the random flicker opacity
    let currentOpacity = $state(0);

    // text sizes
    const textSize = $derived.by(() => {
        if (!trackTitle) return "text-md";

        if (trackTitle.length >= 90) {
            return "text-xs";
        } else if (trackTitle.length > 70) {
            return "text-sm";
        } else if (trackTitle.length > 40) {
            return "text-md";
        } else if (trackTitle.length > 30) {
            return "text-lg";
        } else if (trackTitle.length > 20) {
            return "text-xl";
        } else {
            return "text-3xl";
        }
    });

    // random bell number generator
    // taken from https://stackoverflow.com/a/49434653
    function randn_bm() {
        let u = 0;
        let v = 0;
        while (u === 0) u = Math.random(); //Converting [0,1) to (0,1)
        while (v === 0) v = Math.random();
        let num = Math.sqrt(-2.0 * Math.log(u)) * Math.cos(2.0 * Math.PI * v);

        num = num / 10.0 + 0.5; // Translate to 0 -> 1
        if (num > 1 || num < 0)
            num = randn_bm(); // resample between 0 and 1 if out of range
        else {
            num = Math.pow(num, flickerSkew); // Skew
            num *= flickerOpacityMax - flickerOpacityMin; // Stretch to fill range
            num += flickerOpacityMin; // offset to min
        }
        return num;
    }

    // random changing the opacity of the monitor to create a flickering effect
    onMount(() => {
        const interval = setInterval(() => {
            currentOpacity = randn_bm();
        }, 1000 / flickerFPS);

        return () => clearInterval(interval);
    });

    // crt dimensions
    const crtWidth = 340;
    const crtHeight = 168;
</script>

<div class="relative" style="width: {crtWidth}px; height: {crtHeight}px;">
    <!-- CRT effect -->
    <div class="z-10 absolute w-full h-full crt bg-emerald-900 rounded-lg" style="opacity: {currentOpacity};"></div>
    <!-- <div class="z-10 absolute w-full h-full crt bg-zinc-800 rounded-lg opacity-20"></div> -->
    <div class="z-30 absolute w-full h-full crt bg-zinc-800 rounded-lg opacity-40"></div>

    <!-- Monitor Media -->
    <div class="absolute w-full h-full flex flex-col p-2 text-center">
        <!-- Heading -->
        <span
            class="crt-text text-emerald-200 text-lg h-[15%]"
            style="opacity: {currentOpacity + textFlickerModifier}; font-family: 'PerfectDosVga', sans-serif;"
        >
            {playerStateText}
        </span>

        {#if audio}
            <!-- Track Info -->
            <div class="w-full h-full flex flex-row items-center p-2 text-center">
                <!-- Artist Image & Name -->
                <div class="w-[30%] h-full flex flex-col justify-center">
                    <img
                        src={artistImage}
                        alt=""
                        class="w-[80%] aspect-square object-cover z-20 rounded-lg border-3 border-zinc-900 glitch-box animate-bounce"
                    />

                    <div class="h-2"></div>

                    <span
                        class="crt-text-bounce text-emerald-200 text-sm w-[80%]"
                        style="opacity: {currentOpacity +
                            textFlickerModifier}; font-family: 'PerfectDosVga', sans-serif;"
                    >
                        {artistName}
                    </span>
                </div>

                <!-- Track Name -->
                <span
                    class="crt-text text-emerald-200 w-[40%] {textSize}"
                    style="opacity: {currentOpacity + textFlickerModifier}; font-family: 'PerfectDosVga', sans-serif;"
                >
                    {trackTitle}
                </span>

                <!-- Album Cover & Name -->
                <div class="w-[30%] h-full flex flex-col justify-center">
                    {#if albumCover}
                        <img
                            src={albumCover}
                            alt=""
                            class="w-[80%] aspect-square object-cover z-20 rounded-lg border-3 border-zinc-900 glitch-box-delayed self-end"
                        />

                        <div class="h-2"></div>
                    {/if}

                    <span
                        class="crt-text-bounce-delayed text-emerald-200 text-sm w-[80%] self-end"
                        style="opacity: {currentOpacity +
                            textFlickerModifier}; font-family: 'PerfectDosVga', sans-serif;"
                    >
                        {albumTitle}
                    </span>
                </div>
            </div>
        {/if}
    </div>
</div>

<style>
    /* copied from https://aleclownes.com/2017/02/01/crt-display.html */
    @keyframes textShadow {
        0% {
            text-shadow:
                0.4389924193300864px 0 1px rgba(0, 30, 255, 0.5),
                -0.4389924193300864px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        5% {
            text-shadow:
                2.7928974010788217px 0 1px rgba(0, 30, 255, 0.5),
                -2.7928974010788217px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        10% {
            text-shadow:
                0.02956275843481219px 0 1px rgba(0, 30, 255, 0.5),
                -0.02956275843481219px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        15% {
            text-shadow:
                0.40218538552878136px 0 1px rgba(0, 30, 255, 0.5),
                -0.40218538552878136px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        20% {
            text-shadow:
                3.4794037899852017px 0 1px rgba(0, 30, 255, 0.5),
                -3.4794037899852017px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        25% {
            text-shadow:
                1.6125630401149584px 0 1px rgba(0, 30, 255, 0.5),
                -1.6125630401149584px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        30% {
            text-shadow:
                0.7015590085143956px 0 1px rgba(0, 30, 255, 0.5),
                -0.7015590085143956px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        35% {
            text-shadow:
                3.896914047650351px 0 1px rgba(0, 30, 255, 0.5),
                -3.896914047650351px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        40% {
            text-shadow:
                3.870905614848819px 0 1px rgba(0, 30, 255, 0.5),
                -3.870905614848819px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        45% {
            text-shadow:
                2.231056963361899px 0 1px rgba(0, 30, 255, 0.5),
                -2.231056963361899px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        50% {
            text-shadow:
                0.08084290417898504px 0 1px rgba(0, 30, 255, 0.5),
                -0.08084290417898504px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        55% {
            text-shadow:
                2.3758461067427543px 0 1px rgba(0, 30, 255, 0.5),
                -2.3758461067427543px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        60% {
            text-shadow:
                2.202193051050636px 0 1px rgba(0, 30, 255, 0.5),
                -2.202193051050636px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        65% {
            text-shadow:
                2.8638780614874975px 0 1px rgba(0, 30, 255, 0.5),
                -2.8638780614874975px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        70% {
            text-shadow:
                0.48874025155497314px 0 1px rgba(0, 30, 255, 0.5),
                -0.48874025155497314px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        75% {
            text-shadow:
                1.8948491305757957px 0 1px rgba(0, 30, 255, 0.5),
                -1.8948491305757957px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        80% {
            text-shadow:
                0.0833037308038857px 0 1px rgba(0, 30, 255, 0.5),
                -0.0833037308038857px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        85% {
            text-shadow:
                0.09769827255241735px 0 1px rgba(0, 30, 255, 0.5),
                -0.09769827255241735px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        90% {
            text-shadow:
                3.443339761481782px 0 1px rgba(0, 30, 255, 0.5),
                -3.443339761481782px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        95% {
            text-shadow:
                2.1841838852799786px 0 1px rgba(0, 30, 255, 0.5),
                -2.1841838852799786px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
        100% {
            text-shadow:
                2.6208764473832513px 0 1px rgba(0, 30, 255, 0.5),
                -2.6208764473832513px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px;
        }
    }

    @keyframes boxShadow {
        0% {
            box-shadow:
                0.22px 0 1px rgba(0, 30, 255, 0.1),
                -0.22px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        5% {
            box-shadow:
                1.4px 0 1px rgba(0, 30, 255, 0.1),
                -1.4px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        10% {
            box-shadow:
                0.01px 0 1px rgba(0, 30, 255, 0.1),
                -0.01px 0 1px rgba(255, 0, 80, 0.3),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        15% {
            box-shadow:
                0.2px 0 1px rgba(0, 30, 255, 0.1),
                -0.2px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        20% {
            box-shadow:
                1.7px 0 1px rgba(0, 30, 255, 0.1),
                -1.7px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        25% {
            box-shadow:
                0.8px 0 1px rgba(0, 30, 255, 0.1),
                -0.8px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        30% {
            box-shadow:
                0.35px 0 1px rgba(0, 30, 255, 0.1),
                -0.35px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        35% {
            box-shadow:
                1.9px 0 1px rgba(0, 30, 255, 0.1),
                -1.9px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        40% {
            box-shadow:
                1.9px 0 1px rgba(0, 30, 255, 0.1),
                -1.9px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        45% {
            box-shadow:
                1.1px 0 1px rgba(0, 30, 255, 0.1),
                -1.1px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        50% {
            box-shadow:
                0.04px 0 1px rgba(0, 30, 255, 0.1),
                -0.04px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        55% {
            box-shadow:
                1.2px 0 1px rgba(0, 30, 255, 0.1),
                -1.2px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        60% {
            box-shadow:
                1.1px 0 1px rgba(0, 30, 255, 0.1),
                -1.1px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        65% {
            box-shadow:
                1.4px 0 1px rgba(0, 30, 255, 0.1),
                -1.4px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        70% {
            box-shadow:
                0.24px 0 1px rgba(0, 30, 255, 0.1),
                -0.24px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        75% {
            box-shadow:
                0.95px 0 1px rgba(0, 30, 255, 0.1),
                -0.95px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        80% {
            box-shadow:
                0.04px 0 1px rgba(0, 30, 255, 0.1),
                -0.04px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        85% {
            box-shadow:
                0.05px 0 1px rgba(0, 30, 255, 0.1),
                -0.05px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        90% {
            box-shadow:
                1.7px 0 1px rgba(0, 30, 255, 0.1),
                -1.7px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        95% {
            box-shadow:
                1.1px 0 1px rgba(0, 30, 255, 0.1),
                -1.1px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
        100% {
            box-shadow:
                1.3px 0 1px rgba(0, 30, 255, 0.1),
                -1.3px 0 1px rgba(255, 0, 80, 0.1),
                0 0 3px rgba(0, 0, 0, 0.2);
        }
    }

    @keyframes bounce {
        50% {
            transform: translateY(-10%);
            animation-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
        }
        0%,
        100% {
            transform: none;
            animation-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
        }
    }

    .glitch-box {
        animation:
            boxShadow 1.6s infinite,
            bounce 4s infinite;
    }

    .glitch-box-delayed {
        animation:
            boxShadow 1.6s infinite,
            bounce 3.4s infinite;

        animation-delay: 0.2s;
    }

    .crt::after {
        content: " ";
        display: block;
        position: absolute;
        pointer-events: none;
    }
    .crt::before {
        content: " ";
        position: absolute;
        top: 0;
        left: 0;
        bottom: 0;
        right: 0;
        background:
            linear-gradient(rgba(18, 16, 16, 0) 50%, rgba(0, 0, 0, 0.25) 50%),
            linear-gradient(90deg, rgba(255, 0, 0, 0.06), rgba(0, 255, 0, 0.02), rgba(0, 0, 255, 0.06));
        z-index: 2;
        background-size:
            100% 10px,
            5px 100%;
        pointer-events: none;
    }
    .crt {
        animation: flicker 0.15s infinite;
    }

    .crt-text {
        animation: textShadow 1.6s infinite;
    }

    .crt-text-bounce {
        animation:
            textShadow 1.6s infinite,
            bounce 4s infinite;
    }

    .crt-text-bounce-delayed {
        animation:
            textShadow 1.6s infinite,
            bounce 3.4s infinite;

        animation-delay: 0.2s;
    }
</style>
