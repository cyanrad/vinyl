<script lang="ts">
    import './app.css';
    import Base from './lib/Base.svelte';
    import { setContext } from 'svelte';

    // Setting the FPS context to be used in animations
    setContext('fps', 100);
</script>

<main class="flex items-center justify-center h-screen w-screen fixed inset-0">
    <Base />
</main>

<style>
</style>
