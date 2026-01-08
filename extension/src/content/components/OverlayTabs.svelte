<script lang="ts">
  export type ViewType = "hierarchy" | "dependencies" | "timeline"

  export let activeView: ViewType
  export let onChange: (view: ViewType) => void

  const tabs = [
    { id: "hierarchy", label: "Hierarchy" },
    { id: "dependencies", label: "Dependencies" },
    { id: "timeline", label: "Timeline" }
  ] as const
</script>

<div class="border-b border-[color:var(--pn-border)] px-6 bg-[color:var(--pn-surface-weak)]">
  <div class="flex gap-0">
    {#each tabs as tab}
      <button
        class={`px-4 py-2.5 text-[13.5px] font-medium transition-all relative ${
          activeView === tab.id
            ? "text-[color:var(--pn-text)]"
            : "text-[color:var(--pn-muted-50)] hover:text-[color:var(--pn-text)]"
        }`}
        on:click={() => onChange(tab.id)}
      >
        {tab.label}
        {#if activeView === tab.id}
          <div class="absolute bottom-0 left-0 right-0 h-[2px] bg-[color:var(--pn-text)]" />
        {/if}
      </button>
    {/each}
  </div>
</div>
