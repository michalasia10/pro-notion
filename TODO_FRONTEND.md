# Frontend To-Do List (Browser Extension)

This file outlines the technical tasks required to build the frontend for the Project Overlay Management for Notion application as a browser extension using React and TypeScript.

## Phase 1: Core MVP Setup (4-6 weeks)

### 1. Browser Extension Foundation
- [ ] Set up the basic browser extension structure (Manifest V3 for Chrome/Firefox).
- [ ] Create `content_script` to be injected into `notion.so/*` pages.
- [ ] Create `background_script` for handling API calls and state.
- [ ] Establish communication between content and background scripts (e.g., using `postMessage` or `chrome.runtime.sendMessage`).

### 2. Application Setup
- [ ] Set up a React 18+ project with TypeScript.
- [ ] Integrate a styling solution (e.g., Styled-components).
- [ ] Set up state management with Zustand.
- [ ] Set up an API client for communicating with the backend (e.g., using `fetch` or `axios`).

### 3. Basic Overlay and Notion Integration
- [ ] Implement a `MutationObserver` in the content script to detect when Notion databases are rendered on a page.
- [ ] Inject a floating button or UI element to trigger the overlay.
- [ ] Render a basic, empty overlay component on top of the Notion UI.
- [ ] Implement the frontend part of the OAuth flow to connect to the backend.

## Phase 2: Advanced Features (6-8 weeks)

### 4. Gantt Chart Implementation
- [ ] Choose and integrate a charting library (e.g., D3.js, Vis.js) for rendering the Gantt chart.
- [ ] Fetch task data from the backend and render tasks on the timeline.
- [ ] Implement drag & drop functionality to update task dates.
- [ ] Implement zoom controls (day, week, month views).
- [ ] Display multiple Notion properties (status, assignee) on the task bars.

### 5. Task Dependencies and Hierarchy
- [ ] Visualize task dependencies as lines/arrows between tasks on the Gantt chart.
- [ ] Implement UI for creating and deleting dependencies.
- [ ] Develop a tree-view or nested structure to display hierarchical tasks.
- [ ] Implement collapse/expand functionality for parent tasks.

### 6. Real-time Synchronization
- [ ] Implement a WebSocket client (`socket.io-client`) to listen for real-time updates from the backend.
- [ ] Update the UI (Gantt chart, task lists) in real-time when changes are pushed from the server.

## Phase 3: PM Features (4-6 weeks)

### 7. Advanced UI/UX
- [ ] Design and implement the Multi-Project Dashboard view.
- [ ] Create components for Team Workload Visualization.
- [ ] Develop UI for time tracking (start/stop buttons, timers).
- [ ] Create views for displaying project progress reports.
- [ ] Refine the overall UI/UX, including loading states, error messages, and seamless integration with Notion's design.
