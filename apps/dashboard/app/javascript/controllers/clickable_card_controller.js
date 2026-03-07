import { Controller } from "@hotwired/stimulus";

// Makes an entire card element clickable, navigating to a target URL,
// while still allowing interactive elements (links, buttons) inside the
// card to handle their own click events normally.
//
// Usage:
//   <div
//     data-controller="clickable-card"
//     data-clickable-card-url-value="/examples/my-example"
//     data-action="click->clickable-card#navigate keydown->clickable-card#navigateOnKeydown"
//   >
//     ...
//     <a href="/other-link">This link still works independently</a>
//   </div>
export default class extends Controller {
  static values = {
    url: String,
  };

  navigate(event) {
    if (event.target.closest("a, button")) {
      return;
    }
    window.location.href = this.urlValue;
  }

  navigateOnKeydown(event) {
    if (event.target.closest("a, button")) {
      return;
    }
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      window.location.href = this.urlValue;
    }
  }
}
