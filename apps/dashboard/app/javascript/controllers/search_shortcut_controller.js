import { Controller } from "@hotwired/stimulus";

// Global keyboard shortcuts for search
// Listens for "/" or "Cmd+K"/"Ctrl+K" to focus the navbar search
export default class extends Controller {
  connect() {
    this.handleKeydown = this.handleKeydown.bind(this);
    document.addEventListener("keydown", this.handleKeydown);
  }

  disconnect() {
    document.removeEventListener("keydown", this.handleKeydown);
  }

  handleKeydown(event) {
    // Only trigger if not in an input/textarea
    const inInput = event.target.tagName === "INPUT" ||
      event.target.tagName === "TEXTAREA" ||
      event.target.isContentEditable;

    if (inInput) return;

    // "/" key or "Cmd+K" / "Ctrl+K"
    const isSlashKey = event.key === "/";
    const isCmdK = (event.metaKey || event.ctrlKey) && event.key === "k";

    if (isSlashKey || isCmdK) {
      event.preventDefault();

      // Find and focus navbar search
      const searchElement = document.querySelector(
        '[data-controller~="navbar-search"]',
      );
      if (searchElement) {
        const searchController = this.application
          .getControllerForElementAndIdentifier(
            searchElement,
            "navbar-search",
          );

        if (searchController && typeof searchController.focus === "function") {
          searchController.focus();
        }
      }
    }
  }
}
