export type TaskItem = {
  id: string
  name: string
  status: string
  assignee: string
  type: "project" | "task" | "subtask"
  children?: TaskItem[]
}

export const tasks: TaskItem[] = [
  {
    id: "1",
    name: "Q1 2026 Product Launch",
    status: "In Progress",
    assignee: "Sarah Chen",
    type: "project",
    children: [
      {
        id: "1.1",
        name: "Design System Overhaul",
        status: "In Progress",
        assignee: "Sarah Chen",
        type: "task",
        children: [
          {
            id: "1.1.1",
            name: "Component Library",
            status: "Completed",
            assignee: "Sarah Chen",
            type: "subtask"
          },
          {
            id: "1.1.2",
            name: "Design Tokens",
            status: "In Progress",
            assignee: "Sarah Chen",
            type: "subtask"
          },
          {
            id: "1.1.3",
            name: "Documentation",
            status: "Planning",
            assignee: "Marcus Johnson",
            type: "subtask"
          }
        ]
      },
      {
        id: "1.2",
        name: "API v2 Development",
        status: "Planning",
        assignee: "Marcus Johnson",
        type: "task",
        children: [
          {
            id: "1.2.1",
            name: "Schema Design",
            status: "Planning",
            assignee: "Marcus Johnson",
            type: "subtask"
          },
          {
            id: "1.2.2",
            name: "Implementation",
            status: "Not Started",
            assignee: "Alex Rivera",
            type: "subtask"
          }
        ]
      },
      {
        id: "1.3",
        name: "Marketing Campaign",
        status: "Not Started",
        assignee: "Alex Rivera",
        type: "task"
      }
    ]
  },
  {
    id: "2",
    name: "Q2 2026 Infrastructure",
    status: "Planning",
    assignee: "Marcus Johnson",
    type: "project",
    children: [
      {
        id: "2.1",
        name: "Server Migration",
        status: "Planning",
        assignee: "Marcus Johnson",
        type: "task"
      }
    ]
  }
]

export const dependencies = [
  {
    from: "Design System Overhaul",
    to: "API v2 Development",
    type: "Finish-to-Start",
    lag: "3 days"
  },
  {
    from: "API v2 Development",
    to: "Mobile App Beta",
    type: "Finish-to-Start",
    lag: "0 days"
  },
  {
    from: "Mobile App Beta",
    to: "Testing & QA",
    type: "Finish-to-Start",
    lag: "1 day"
  }
]

export const ganttTasks = [
  {
    name: "Design System Overhaul",
    assignee: "Sarah Chen",
    status: "In Progress",
    start: 10,
    duration: 45,
    color: "bg-blue-500"
  },
  {
    name: "API v2 Development",
    assignee: "Marcus Johnson",
    status: "Planning",
    start: 30,
    duration: 60,
    color: "bg-emerald-500"
  },
  {
    name: "Mobile App Beta",
    assignee: "Emma Wilson",
    status: "Planning",
    start: 55,
    duration: 50,
    color: "bg-purple-500"
  },
  {
    name: "Testing & QA",
    assignee: "Alex Rivera",
    status: "Not Started",
    start: 70,
    duration: 35,
    color: "bg-orange-500"
  }
]

export const months = ["January", "February", "March", "April"]
