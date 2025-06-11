<script lang="ts">
    import { setContext } from "svelte";

    // styles
    import "./app.css";

    // components
    import Base from "./lib/Base.svelte";

    // api
    import { API_URL } from "./lib/api/consts";
    import PocketBase from "pocketbase";

    // Setting the FPS context to be used in animations
    setContext("fps", 100);
    setContext("pb", new PocketBase(API_URL));
</script>

<main class="flex items-center justify-center h-screen w-screen fixed inset-0">
    <Base />
</main>

<style>
</style>
