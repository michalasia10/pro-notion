<script lang="ts">
  export let onToggle: () => void;
</script>

<div
  style="position:fixed;right:24px;bottom:24px;z-index:2147483647;pointer-events:auto;display:block;visibility:visible;"
>
  <button
    class="group flex items-center gap-2.5 rounded-lg bg-[color:var(--pn-button-bg)] px-4 py-2.5 text-[13.5px] font-medium text-[color:var(--pn-button-text)] shadow-lg transition-all hover:shadow-xl hover:scale-[1.02]"
    on:click={onToggle}
  >
    <svg
      class="h-4 w-4 transition-transform group-hover:scale-110"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      stroke-width="1.5"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"
      />
    </svg>
    Pro Notion
  </button>
</div>
