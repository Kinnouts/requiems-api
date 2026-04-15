import { Controller } from "@hotwired/stimulus";

// Connects to data-controller="email-confirm"
//
// A generic controller for destructive actions requiring typed confirmation.
// The submit button stays disabled until the user types the exact expected string.
//
// Values:
//   expected (String) — the exact string the user must type to unlock submit
//
// Targets:
//   modal  — the full-screen backdrop overlay
//   input  — the confirmation text input
//   submit — the form submit button
//   hint   — (optional) element for live match feedback
//
// Actions (wire via data-action):
//   open            — show the modal            (on trigger button)
//   close           — hide the modal            (on cancel/close buttons)
//   validate        — re-check on each keystroke (on input, event: input)
//   closeOnBackdrop — close when clicking backdrop (on modal div, event: click)
//
// Example:
//   <div data-controller="email-confirm"
//        data-email-confirm-expected-value="<%= current_user.email %>">
//     <button data-action="click->email-confirm#open">Delete</button>
//     <div data-email-confirm-target="modal"
//          data-action="click->email-confirm#closeOnBackdrop"
//          class="hidden fixed inset-0 z-50 ...">
//       ...
//       <input data-email-confirm-target="input"
//              data-action="input->email-confirm#validate" />
//       <p data-email-confirm-target="hint"></p>
//       <button data-email-confirm-target="submit" disabled>Confirm</button>
//     </div>
//   </div>
export default class extends Controller {
  static targets = ["modal", "input", "submit", "hint"];
  static values = { expected: String };

  connect() {
    this._escapeHandler = this._onEscape.bind(this);
    document.addEventListener("keydown", this._escapeHandler);
  }

  disconnect() {
    document.removeEventListener("keydown", this._escapeHandler);
    document.body.classList.remove("overflow-hidden");
  }

  open(event) {
    event?.preventDefault();
    this.modalTarget.classList.remove("hidden");
    document.body.classList.add("overflow-hidden");
    requestAnimationFrame(() => this.inputTarget?.focus());
  }

  close(event) {
    event?.preventDefault();
    this.modalTarget.classList.add("hidden");
    document.body.classList.remove("overflow-hidden");
    this._reset();
  }

  closeOnBackdrop(event) {
    if (event.target === this.modalTarget) this.close(event);
  }

  validate() {
    const value = this.inputTarget.value;
    const match = value === this.expectedValue;

    this.submitTarget.disabled = !match;

    // Input border feedback
    this.inputTarget.classList.toggle("border-green-500", match);
    this.inputTarget.classList.toggle("dark:border-green-400", match);
    this.inputTarget.classList.toggle("ring-2", match);
    this.inputTarget.classList.toggle("ring-green-500", match);
    this.inputTarget.classList.toggle("border-gray-300", !match);
    this.inputTarget.classList.toggle("dark:border-gray-600", !match);

    // Hint text
    if (!this.hasHintTarget) return;

    if (match) {
      this.hintTarget.textContent = "✓ Confirmed — you may now proceed";
      this.hintTarget.className =
        "text-xs font-medium text-green-600 dark:text-green-400 mt-1.5";
    } else if (value.length > 0) {
      this.hintTarget.textContent =
        "Doesn't match — type your exact email address";
      this.hintTarget.className =
        "text-xs text-red-500 dark:text-red-400 mt-1.5";
    } else {
      this.hintTarget.textContent = "";
      this.hintTarget.className = "text-xs mt-1.5";
    }
  }

  // — Private —

  _onEscape(event) {
    if (
      event.key === "Escape" && !this.modalTarget.classList.contains("hidden")
    ) {
      this.close();
    }
  }

  _reset() {
    if (this.hasInputTarget) {
      this.inputTarget.value = "";
      this.inputTarget.classList.remove(
        "border-green-500",
        "dark:border-green-400",
        "ring-2",
        "ring-green-500",
      );
      this.inputTarget.classList.add("border-gray-300", "dark:border-gray-600");
    }

    if (this.hasSubmitTarget) {
      this.submitTarget.disabled = true;
    }

    if (this.hasHintTarget) {
      this.hintTarget.textContent = "";
      this.hintTarget.className = "text-xs mt-1.5";
    }
  }
}
