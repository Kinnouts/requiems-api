import { Controller } from "@hotwired/stimulus";

// Swaps the displayed price between monthly and yearly on the private deployment form.
export default class extends Controller {
  static targets = ["monthlyPrice", "yearlyPrice"];

  connect() {
    const selected = this.element.querySelector(
      "input[name$='[billing_cycle]']:checked",
    );
    this.#applyVisibility(selected?.value === "yearly");
  }

  toggle(event) {
    this.#applyVisibility(event.target.value === "yearly");
  }

  #applyVisibility(isYearly) {
    this.monthlyPriceTargets.forEach((el) =>
      el.classList.toggle("hidden", isYearly)
    );
    this.yearlyPriceTargets.forEach((el) =>
      el.classList.toggle("hidden", !isYearly)
    );
  }
}
