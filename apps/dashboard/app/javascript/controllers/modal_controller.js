import { Controller } from "@hotwired/stimulus";

// Modal dialog controller
// Usage: data-controller="modal"
export default class extends Controller {
  open() {
    this.element.classList.remove("hidden");
    document.body.classList.add("overflow-hidden");
  }

  close() {
    this.element.classList.add("hidden");
    document.body.classList.remove("overflow-hidden");
  }

  // Close on Escape key
  closeWithKeyboard(event) {
    if (event.key === "Escape") {
      this.close();
    }
  }

  connect() {
    // Listen for Escape key
    document.addEventListener("keydown", this.closeWithKeyboard.bind(this));
  }

  disconnect() {
    document.removeEventListener("keydown", this.closeWithKeyboard.bind(this));
    document.body.classList.remove("overflow-hidden");
  }
}
