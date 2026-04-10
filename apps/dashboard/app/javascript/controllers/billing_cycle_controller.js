import { Controller } from "@hotwired/stimulus"

// Swaps the displayed price between monthly and yearly on the private deployment form.
export default class extends Controller {
  static targets = ["monthlyPrice", "yearlyPrice"]

  toggle(event) {
    const isYearly = event.target.value === "yearly"

    this.monthlyPriceTargets.forEach(el => el.classList.toggle("hidden", isYearly))
    this.yearlyPriceTargets.forEach(el => el.classList.toggle("hidden", !isYearly))
  }
}
