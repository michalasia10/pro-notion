<script lang="ts">
  import { tasks } from "./data"
  import type { TaskItem } from "./data"

  type FlattenedTask = {
    task: TaskItem
    level: number
    hasChildren: boolean
    canExpand: boolean
  }

  let maxDepth = 3

  function seedExpanded(items: TaskItem[], level = 0, set = new Set<string>()) {
    for (const item of items) {
      if (level < 2) {
        set.add(item.id)
      }
      if (item.children) {
        seedExpanded(item.children, level + 1, set)
      }
    }
    return set
  }

  let expandedIds = seedExpanded(tasks)

  function toggleExpand(id: string) {
    const next = new Set(expandedIds)
    if (next.has(id)) {
      next.delete(id)
    } else {
      next.add(id)
    }
    expandedIds = next
  }

  function buildTaskList(items: TaskItem[], level = 0, acc: FlattenedTask[] = []) {
    for (const item of items) {
      const hasChildren = Boolean(item.children && item.children.length > 0)
      const canExpand = hasChildren && level + 1 < maxDepth
      acc.push({ task: item, level, hasChildren, canExpand })
      if (hasChildren && canExpand && expandedIds.has(item.id)) {
        buildTaskList(item.children || [], level + 1, acc)
      }
    }
    return acc
  }

  $: flattenedTasks = buildTaskList(tasks)
</script>

<div class="h-full flex flex-col max-w-5xl mx-auto">
  <div class="flex items-center justify-between pb-4 mb-4 border-b border-[color:var(--pn-border)]">
    <div>
      <h3 class="text-[15px] font-semibold text-[color:var(--pn-text)] mb-1">
        Project Hierarchy
      </h3>
      <p class="text-[12px] text-[color:var(--pn-muted-50)]">
        Showing all projects, tasks, and subtasks
      </p>
    </div>

    <div class="flex items-center gap-2">
      <span class="text-[12px] text-[color:var(--pn-muted)]">Depth:</span>
      <select
        bind:value={maxDepth}
        class="text-[13px] px-2.5 py-1 rounded-md border border-[color:var(--pn-border-strong)] bg-[color:var(--pn-panel)] text-[color:var(--pn-text)] focus:outline-none focus:ring-1 focus:ring-[color:var(--pn-muted-40)]"
      >
        <option value={1}>Level 1 (Projects)</option>
        <option value={2}>Level 2 (+ Tasks)</option>
        <option value={3}>Level 3 (+ Subtasks)</option>
        <option value={99}>All levels</option>
      </select>
    </div>
  </div>

  <div class="flex-1 overflow-y-auto -mx-2">
    <div class="space-y-0.5">
      {#each flattenedTasks as row (row.task.id)}
        <div
          class="group flex items-center gap-2.5 py-2 px-2 hover:bg-[color:var(--pn-hover)] rounded-md transition-colors cursor-pointer"
          style={`padding-left:${row.level * 20 + 8}px`}
        >
          {#if row.hasChildren && row.canExpand}
            <button
              class="flex h-5 w-5 shrink-0 items-center justify-center text-[color:var(--pn-muted-40)] hover:text-[color:var(--pn-text)] hover:bg-[color:var(--pn-hover-strong)] rounded transition-all"
              on:click={() => toggleExpand(row.task.id)}
            >
              <svg
                class={`h-3.5 w-3.5 transition-transform duration-150 ${
                  expandedIds.has(row.task.id) ? "rotate-90" : ""
                }`}
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                stroke-width="2"
              >
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          {:else}
            <div class="w-5 shrink-0" />
          {/if}

          <div
            class={`shrink-0 ${
              row.task.type === "project"
                ? "h-5 w-5 text-[color:var(--pn-text)]"
                : row.task.type === "task"
                  ? "h-4 w-4 text-[color:var(--pn-muted)]"
                  : "h-3.5 w-3.5 text-[color:var(--pn-muted-50)]"
            }`}
          >
            {#if row.task.type === "project"}
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"
                />
              </svg>
            {:else if row.task.type === "task"}
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12h3.75M9 15h3.75M9 18h3.75m3 .75H18a2.25 2.25 0 002.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 00-1.123-.08m-5.801 0c-.065.21-.1.433-.1.664 0 .414.336.75.75.75h4.5a.75.75 0 00.75-.75 2.25 2.25 0 00-.1-.664m-5.8 0A2.251 2.251 0 0113.5 2.25H15c1.012 0 1.867.668 2.15 1.586m-5.8 0c-.376.023-.75.05-1.124.08C9.095 4.01 8.25 4.973 8.25 6.108V8.25m0 0H4.875c-.621 0-1.125.504-1.125 1.125v11.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125V9.375c0-.621-.504-1.125-1.125-1.125H8.25zM6.75 12h.008v.008H6.75V12zm0 3h.008v.008H6.75V15zm0 3h.008v.008H6.75V18z"
                />
              </svg>
            {:else}
              <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            {/if}
          </div>

          <span class="flex-1 text-[13.5px] text-[color:var(--pn-text)] font-normal min-w-0 truncate">
            {row.task.name}
          </span>

          <div class="flex items-center gap-2 shrink-0">
            <div class="flex items-center gap-1.5">
              <div
                class={`h-1.5 w-1.5 rounded-full ${
                  row.task.status === "Completed"
                    ? "bg-emerald-500"
                    : row.task.status === "In Progress"
                      ? "bg-blue-500"
                      : row.task.status === "Planning"
                        ? "bg-amber-500"
                        : "bg-slate-300"
                }`}
              />
              <span class="text-[12px] text-[color:var(--pn-muted)]">{row.task.status}</span>
            </div>

            <div
              class="h-5 w-5 rounded-full bg-gradient-to-br from-blue-400 to-indigo-500 flex items-center justify-center text-white text-[9px] font-medium"
            >
              {row.task.assignee
                .split(" ")
                .map((name) => name[0])
                .join("")}
            </div>
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>
