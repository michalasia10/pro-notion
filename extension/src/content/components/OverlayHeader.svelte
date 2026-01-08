<script lang="ts">
  export type StatusConfig = {
    label: string
    color: string
    dot: string
  }

  export let currentStatus: StatusConfig
  export let onClose: () => void
</script>

<div class="flex items-center justify-between border-b border-[color:var(--pn-border)] px-6 py-3.5 bg-[color:var(--pn-surface)]">
  <div class="flex items-center gap-4">
    <div class="flex flex-col">
      <div class="flex items-center gap-2">
        <h2 class="text-[15px] font-semibold text-[color:var(--pn-text)]">Project Alpha</h2>
        <span class="text-[13px] text-[color:var(--pn-muted-40)]">/</span>
        <span class="text-[13px] text-[color:var(--pn-muted)]">Tasks Database</span>
      </div>
    </div>

    <div
      class={`flex items-center gap-1.5 px-2.5 py-1 rounded-full border text-[12px] font-medium ${currentStatus.color}`}
    >
      <div class={`h-1.5 w-1.5 rounded-full ${currentStatus.dot}`} />
      {currentStatus.label}
    </div>
  </div>

  <button
    class="flex h-7 w-7 items-center justify-center rounded-md text-[color:var(--pn-muted)] hover:bg-[color:var(--pn-hover-strong)]"
    on:click={onClose}
  >
    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
      <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
    </svg>
  </button>
</div>
