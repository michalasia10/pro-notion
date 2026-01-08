<script lang="ts">
  import { ganttTasks, months } from "./data"
</script>

<div class="space-y-6 max-w-5xl mx-auto">
  <div class="flex items-center justify-between">
    <div>
      <h3 class="text-[15px] font-semibold text-[color:var(--pn-text)]">Q1 2026 Timeline</h3>
      <p class="text-[12px] text-[color:var(--pn-muted-50)]">
        Drag tasks to adjust dates and dependencies
      </p>
    </div>
    <div class="flex gap-2">
      <button class="rounded-md border border-[color:var(--pn-border-strong)] bg-[color:var(--pn-panel)] px-3 py-1.5 text-sm hover:bg-[color:var(--pn-hover)]">
        Today
      </button>
      <button class="rounded-md border border-[color:var(--pn-border-strong)] bg-[color:var(--pn-panel)] px-3 py-1.5 text-sm hover:bg-[color:var(--pn-hover)]">
        Zoom
      </button>
    </div>
  </div>

  <div class="overflow-x-auto rounded-lg border border-[color:var(--pn-border)] bg-[color:var(--pn-panel)]">
    <div class="flex border-b border-[color:var(--pn-border)]">
      <div class="w-80 shrink-0 border-r border-[color:var(--pn-border)] bg-[color:var(--pn-muted-bg)] px-4 py-3 text-sm font-medium">
        Task
      </div>
      <div class="flex flex-1">
        {#each months as month}
          <div
            class="flex-1 border-r border-[color:var(--pn-border)] bg-[color:var(--pn-muted-bg)] px-4 py-3 text-center text-sm font-medium last:border-r-0"
          >
            {month}
          </div>
        {/each}
      </div>
    </div>

    {#each ganttTasks as item, index}
      <div class="flex border-b border-[color:var(--pn-border)] last:border-b-0">
        <div class="w-80 shrink-0 space-y-1 border-r border-[color:var(--pn-border)] px-4 py-3">
          <div class="font-medium text-[color:var(--pn-text)]">{item.name}</div>
          <div class="flex items-center gap-2 text-xs text-[color:var(--pn-muted-50)]">
            <span class="inline-flex items-center gap-1">
              <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                  clip-rule="evenodd"
                />
              </svg>
              {item.assignee}
            </span>
            <span>•</span>
            <span class="rounded-full bg-[color:var(--pn-chip-bg)] px-2 py-0.5">
              {item.status}
            </span>
          </div>
        </div>
        <div class="relative flex flex-1 items-center px-2 py-3">
          <div
            class={`${item.color} absolute h-8 cursor-pointer rounded transition-all hover:opacity-80`}
            style={`left:${item.start}%;width:${item.duration}%`}
          />
        </div>
      </div>
    {/each}
  </div>
</div>
