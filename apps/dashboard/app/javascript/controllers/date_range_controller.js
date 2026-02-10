import { Controller } from "@hotwired/stimulus"

// Connects to data-controller="date-range"
export default class extends Controller {
  static targets = ["form"]

  toggle(event) {
    event.preventDefault()

    if (this.hasFormTarget) {
      this.formTarget.classList.toggle("hidden")
    }
  }
}
