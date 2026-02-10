import { Controller } from "@hotwired/stimulus"

// Client-side API search and filtering
export default class extends Controller {
  static targets = ["input", "card", "grid", "empty", "count", "resultCount"]

  connect() {
    this.totalCount = this.cardTargets.length
  }

  filter() {
    const query = this.inputTarget.value.toLowerCase().trim()
    let visibleCount = 0

    this.cardTargets.forEach((card) => {
      const name = card.dataset.apiName || ""
      const description = card.dataset.apiDescription || ""
      const tags = card.dataset.apiTags || ""

      const matches = name.includes(query) ||
                     description.includes(query) ||
                     tags.includes(query)

      if (query === "" || matches) {
        card.classList.remove("hidden")
        visibleCount++
      } else {
        card.classList.add("hidden")
      }
    })

    // Update result count
    if (this.hasResultCountTarget) {
      this.resultCountTarget.textContent = visibleCount
    }

    // Show/hide empty state
    if (this.hasEmptyTarget && this.hasGridTarget) {
      if (visibleCount === 0 && query !== "") {
        this.gridTarget.classList.add("hidden")
        this.emptyTarget.classList.remove("hidden")
      } else {
        this.gridTarget.classList.remove("hidden")
        this.emptyTarget.classList.add("hidden")
      }
    }
  }
}
