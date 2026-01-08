import cssText from "../app.css?inline"
import Overlay from "./Overlay.svelte"

const rootId = "pro-notion-extension-root"
const styleId = "pro-notion-extension-style"

if (!document.getElementById(rootId)) {
  if (!document.getElementById(styleId)) {
    const style = document.createElement("style")
    style.id = styleId
    style.textContent = cssText
    document.documentElement.appendChild(style)
  }

  const container = document.createElement("div")
  container.id = rootId
  document.body.appendChild(container)

  new Overlay({
    target: container
  })
}
