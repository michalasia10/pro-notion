<script lang="ts">
  export let onToggle: () => void
</script>

<button
  class="fixed bottom-4 left-4 text-[11px] px-2 py-1 rounded bg-[color:var(--pn-hover)] text-[color:var(--pn-muted-40)] hover:bg-[color:var(--pn-hover-strong)]"
  on:click={onToggle}
>
  Toggle Status (Demo)
</button>
