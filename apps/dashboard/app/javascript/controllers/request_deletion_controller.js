import { Controller } from "@hotwired/stimulus";

// Connects to data-controller="request-deletion"
//
// Handles the "Delete Account" danger zone modal.
// The submit button is disabled until the user types a reason of at least
// MIN_REASON_LENGTH characters, so accidental submissions are prevented.
//
// Targets:
//   modal    — the full-screen backdrop overlay
//   textarea — the deletion reason textarea
//   submit   — the form submit button
//   counter  — (optional) character count / hint element
//
// Actions (wire via data-action):
//   open            — show the modal            (on trigger button click)
//   close           — hide the modal            (on cancel/close buttons)
//   validate        — re-check on each keystroke (on textarea, event: input)
//   closeOnBackdrop — close when clicking backdrop (on modal div, event: click)
const MIN_REASON_LENGTH = 10;

export default class extends Controller {
  static targets = ["modal", "textarea", "submit", "counter"];

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
    requestAnimationFrame(() => this.textareaTarget?.focus());
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
    const length = this.textareaTarget.value.trim().length;
    const ready = length >= MIN_REASON_LENGTH;

    this.submitTarget.disabled = !ready;

    if (this.hasCounterTarget) {
      if (length === 0) {
        this.counterTarget.textContent = "";
        this.counterTarget.className = "text-xs mt-1.5";
      } else if (!ready) {
        const remaining = MIN_REASON_LENGTH - length;
        this.counterTarget.textContent = `${remaining} more character${
          remaining === 1 ? "" : "s"
        } needed`;
        this.counterTarget.className =
          "text-xs text-amber-600 dark:text-amber-400 mt-1.5";
      } else {
        this.counterTarget.textContent = `${length} characters`;
        this.counterTarget.className =
          "text-xs text-green-600 dark:text-green-400 mt-1.5";
      }
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
    if (this.hasTextareaTarget) this.textareaTarget.value = "";
    if (this.hasSubmitTarget) this.submitTarget.disabled = true;
    if (this.hasCounterTarget) {
      this.counterTarget.textContent = "";
      this.counterTarget.className = "text-xs mt-1.5";
    }
  }
}
