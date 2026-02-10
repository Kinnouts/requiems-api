import { Controller } from "@hotwired/stimulus"

// Connects to data-controller="confirm-delete"
export default class extends Controller {
  static targets = ["modal", "emailInput"]

  open(event) {
    event.preventDefault()

    if (this.hasModalTarget) {
      this.modalTarget.classList.remove("hidden")
      document.body.style.overflow = "hidden"

      // Focus on email input
      if (this.hasEmailInputTarget) {
        setTimeout(() => {
          this.emailInputTarget.focus()
        }, 100)
      }
    }
  }

  close(event) {
    if (event) {
      event.preventDefault()
    }

    if (this.hasModalTarget) {
      this.modalTarget.classList.add("hidden")
      document.body.style.overflow = ""

      // Clear email input
      if (this.hasEmailInputTarget) {
        this.emailInputTarget.value = ""
      }
    }
  }

  // Close modal when clicking outside
  clickOutside(event) {
    if (event.target === this.modalTarget) {
      this.close()
    }
  }

  // Close modal on Escape key
  handleEscape(event) {
    if (event.key === "Escape") {
      this.close()
    }
  }

  connect() {
    // Add event listener for Escape key
    this.boundHandleEscape = this.handleEscape.bind(this)
    document.addEventListener("keydown", this.boundHandleEscape)

    // Add event listener for clicking outside
    if (this.hasModalTarget) {
      this.modalTarget.addEventListener("click", this.clickOutside.bind(this))
    }
  }

  disconnect() {
    // Remove event listeners
    document.removeEventListener("keydown", this.boundHandleEscape)
  }
}
