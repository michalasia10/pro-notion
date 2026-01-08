<script lang="ts">
  import { dependencies } from "./data"
</script>

<div class="space-y-6 max-w-4xl mx-auto">
  <div>
    <h3 class="text-[15px] font-semibold text-[color:var(--pn-text)]">Task Dependencies</h3>
    <p class="text-[12px] text-[color:var(--pn-muted-50)]">
      Manage dependencies between tasks and automatically reschedule
    </p>
  </div>

  <div class="space-y-4">
    {#each dependencies as dep, index}
      <div
        class="group relative overflow-hidden rounded-lg border border-[color:var(--pn-border)] bg-[color:var(--pn-panel)] p-5 transition-all hover:border-[color:var(--pn-border-strong)] hover:shadow-md"
      >
        <div class="flex items-center justify-between">
          <div class="flex-1 space-y-3">
            <div class="flex items-center gap-3">
              <div
                class="flex h-8 w-8 items-center justify-center rounded-full bg-[color:var(--pn-chip-bg)] text-sm font-semibold text-[color:var(--pn-text)]"
              >
                {index + 1}
              </div>
              <div class="flex-1">
                <div class="flex items-center gap-2">
                  <span class="rounded bg-[color:var(--pn-chip-bg)] px-2 py-1 text-sm font-medium text-[color:var(--pn-text)]">
                    {dep.from}
                  </span>
                  <svg
                    class="h-4 w-4 text-[color:var(--pn-muted-40)]"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M13 7l5 5m0 0l-5 5m5-5H6"
                    />
                  </svg>
                  <span class="rounded bg-[color:var(--pn-chip-bg)] px-2 py-1 text-sm font-medium text-[color:var(--pn-text)]">
                    {dep.to}
                  </span>
                </div>
              </div>
            </div>
            <div class="ml-11 flex items-center gap-4 text-sm text-[color:var(--pn-muted-50)]">
              <span class="flex items-center gap-1">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M13 10V3L4 14h7v7l9-11h-7z"
                  />
                </svg>
                {dep.type}
              </span>
              <span>•</span>
              <span>Lag: {dep.lag}</span>
            </div>
          </div>
          <button class="opacity-0 transition-opacity group-hover:opacity-100">
            <svg
              class="h-5 w-5 text-[color:var(--pn-muted-50)] hover:text-[color:var(--pn-text)]"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"
              />
            </svg>
          </button>
        </div>
      </div>
    {/each}
  </div>

  <button
    class="flex w-full items-center justify-center gap-2 rounded-lg border-2 border-dashed border-[color:var(--pn-border-soft)] py-3 text-sm font-medium text-[color:var(--pn-muted)] transition-colors hover:border-[color:var(--pn-border-strong)] hover:text-[color:var(--pn-text)]"
  >
    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
    </svg>
    Add New Dependency
  </button>
</div>
