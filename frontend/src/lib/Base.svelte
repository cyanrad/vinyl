<script>
  import { onMount } from 'svelte';
  import Vinyl from './Vinyl.svelte';

  let baseUrl = '';

  onMount(async () => {
    const module = await import('../assets/base.png');
    baseUrl = module.default;
  });
</script>

<div class="relative">
  {#if baseUrl}
    <img src={baseUrl} class="h-[620px] w-auto" />
  {/if}
  <Vinyl status="stopped" coverUrl="../assets/phantasmagoria.jpg" />
</div>
