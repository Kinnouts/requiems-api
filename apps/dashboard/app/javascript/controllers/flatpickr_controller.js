import { Controller } from "@hotwired/stimulus"
import flatpickr from "flatpickr"

// Connects to data-controller="flatpickr"
export default class extends Controller {
  static values = {
    mode: { type: String, default: "single" }, // single, range, multiple
    enableTime: { type: Boolean, default: false },
    dateFormat: { type: String, default: "Y-m-d" },
    minDate: String,
    maxDate: String,
    defaultDate: String
  }

  connect() {
    const options = {
      mode: this.modeValue,
      enableTime: this.enableTimeValue,
      dateFormat: this.dateFormatValue
    }

    if (this.hasMinDateValue) {
      options.minDate = this.minDateValue
    }

    if (this.hasMaxDateValue) {
      options.maxDate = this.maxDateValue
    }

    if (this.hasDefaultDateValue) {
      options.defaultDate = this.defaultDateValue
    }

    // Initialize flatpickr on this element
    this.picker = flatpickr(this.element, options)
  }

  disconnect() {
    if (this.picker) {
      this.picker.destroy()
    }
  }
}
