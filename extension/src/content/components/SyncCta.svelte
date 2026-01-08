<script lang="ts">
  export let onConnect: () => void
</script>

<div class="flex h-full items-center justify-center">
  <div class="flex max-w-sm flex-col items-center gap-4 text-center">
    <div class="h-12 w-12 rounded-lg bg-amber-100 flex items-center justify-center">
      <svg
        class="h-6 w-6 text-amber-600"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        stroke-width="1.5"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M13.5 10.5V6.75a4.5 4.5 0 119 0v3.75M3.75 21.75h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H3.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z"
        />
      </svg>
    </div>
    <div>
      <h3 class="text-[15px] font-semibold text-[color:var(--pn-text)] mb-1.5">
        Connect your Notion workspace
      </h3>
      <p class="text-[13px] text-[color:var(--pn-muted)]">
        Connect and sync to view your project hierarchy and dependencies
      </p>
    </div>
    <button
      class="mt-2 bg-[color:var(--pn-button-bg)] text-[color:var(--pn-button-text)] hover:opacity-90 shadow-sm px-4 py-2 rounded-md text-[13px] font-medium"
      on:click={onConnect}
    >
      Connect &amp; Sync
    </button>
  </div>
</div>
