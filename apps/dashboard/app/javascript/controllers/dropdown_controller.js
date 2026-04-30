import { Controller } from "@hotwired/stimulus";

// Dropdown menu controller
// Usage:
//   <div data-controller="dropdown">
//     <button data-action="click->dropdown#toggle">Menu</button>
//     <div data-dropdown-target="menu" class="hidden">
//       <!-- Dropdown content -->
//     </div>
//   </div>
export default class extends Controller {
  static targets = ["menu"];

  toggle(event) {
    event.stopPropagation();
    this.menuTarget.classList.toggle("hidden");
  }

  hide(event) {
    if (!this.element.contains(event.target)) {
      this.menuTarget.classList.add("hidden");
    }
  }

  connect() {
    // Close dropdown when clicking outside
    document.addEventListener("click", this.hide.bind(this));
  }

  disconnect() {
    document.removeEventListener("click", this.hide.bind(this));
  }
}
