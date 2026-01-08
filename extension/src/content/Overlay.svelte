<script lang="ts">
  import { onMount } from "svelte"
  import {
    DemoStatusToggle,
    DependenciesView,
    FloatingButton,
    HierarchyView,
    OverlayHeader,
    OverlayTabs,
    SyncCta,
    TimelineView
  } from "./components"
  import { getCurrentUser } from "./api"

  type ViewType = "hierarchy" | "dependencies" | "timeline"
  type SyncStatus = "not-connected" | "syncing" | "up-to-date" | "error"

  let isOpen = false
  let activeView: ViewType = "hierarchy"
  let syncStatus: SyncStatus = "syncing"
  let userName = "Project Alpha"
  let workspaceName = "Tasks Database"

  const statusConfig = {
    "not-connected": {
      label: "Not connected",
      color: "bg-amber-50 text-amber-700 border-amber-200",
      dot: "bg-amber-500"
    },
    syncing: {
      label: "Syncing...",
      color: "bg-blue-50 text-blue-700 border-blue-200",
      dot: "bg-blue-500 animate-pulse"
    },
    "up-to-date": {
      label: "Up to date",
      color: "bg-emerald-50 text-emerald-700 border-emerald-200",
      dot: "bg-emerald-500"
    },
    error: {
      label: "Sync error",
      color: "bg-red-50 text-red-700 border-red-200",
      dot: "bg-red-500"
    }
  }

  function toggle() {
    isOpen = !isOpen
  }

  function close() {
    isOpen = false
  }

  async function connectNotion() {
    syncStatus = "syncing"
    try {
      const result = await new Promise<{ ok: boolean; error?: string }>((resolve) => {
        chrome.runtime.sendMessage({ type: "pn-auth-start" }, resolve)
      })
      if (!result.ok) {
        syncStatus = "error"
        return
      }
      await refreshCurrentUser()
    } catch (_error) {
      syncStatus = "error"
    }
  }

  async function refreshCurrentUser() {
    syncStatus = "syncing"
    const user = await getCurrentUser()
    if (!user) {
      syncStatus = "not-connected"
      return
    }

    userName = user.name || "Project Alpha"
    workspaceName = user.notion_workspace_id ? "Workspace" : "Tasks Database"
    syncStatus = user.has_notion_integration ? "up-to-date" : "not-connected"
  }

  function cycleStatus() {
    const statuses: SyncStatus[] = ["not-connected", "syncing", "up-to-date", "error"]
    const currentIndex = statuses.indexOf(syncStatus)
    syncStatus = statuses[(currentIndex + 1) % statuses.length]
  }

  onMount(() => {
    const handleEsc = (event: KeyboardEvent) => {
      if (event.key === "Escape" && isOpen) {
        close()
      }
    }
    window.addEventListener("keydown", handleEsc)
    refreshCurrentUser()
    return () => window.removeEventListener("keydown", handleEsc)
  })

  $: currentStatus = statusConfig[syncStatus]
</script>

<div class="pn-overlay-root" style="position:relative;z-index:2147483647;">
  <FloatingButton onToggle={toggle} />

  {#if isOpen}
    <div class="fixed inset-0 z-[2147483646] bg-black/20 backdrop-blur-[2px]">
      <div class="relative h-screen w-screen bg-[color:var(--pn-bg)] shadow-xl">
        <OverlayHeader
          currentStatus={currentStatus}
          onClose={close}
          title={userName}
          subtitle={workspaceName}
        />
        <OverlayTabs activeView={activeView} onChange={(view) => (activeView = view)} />

        <div class="h-[calc(100vh-110px)] overflow-auto px-6 py-5">
          {#if syncStatus === "not-connected"}
            <SyncCta onConnect={connectNotion} />
          {:else}
            {#if activeView === "hierarchy"}
              <HierarchyView />
            {:else if activeView === "dependencies"}
              <DependenciesView />
            {:else}
              <TimelineView />
            {/if}
          {/if}
        </div>

        <DemoStatusToggle onToggle={cycleStatus} />
      </div>
    </div>
  {/if}
</div>
