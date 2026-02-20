import { Controller } from "@hotwired/stimulus"

// Interactive API playground for testing endpoints — routes through /api/proxy
export default class extends Controller {
  static targets = [
    "param",
    "submitButton",
    "submitText",
    "loading",
    "responseContainer",
    "responseStatus",
    "statusBadge",
    "responseTime",
    "responseBody",
    "errorContainer",
    "errorMessage"
  ]

  static values = {
    method: String,
    url: String,
    index: Number
  }

  connect() {
    this.startTime = null
  }

  async sendRequest(event) {
    event.preventDefault()

    // Collect parameters
    const requestData = this.collectParameters()

    // Validate required parameters
    if (!this.validateParameters()) {
      this.showError("Please fill in all required fields")
      return
    }

    // Show loading state
    this.showLoading()

    try {
      this.startTime = performance.now()

      const response = await this.makeRequest(requestData)

      const endTime = performance.now()
      const clientDuration = Math.round(endTime - this.startTime)

      await this.handleResponse(response, clientDuration)
    } catch (error) {
      this.showError(error.message || "Request failed. Please try again.")
    } finally {
      this.hideLoading()
    }
  }

  collectParameters() {
    const data = {}

    this.paramTargets.forEach((input) => {
      const name = input.dataset.paramName
      const type = input.dataset.paramType
      let value = input.value.trim()

      if (!value) return

      // Parse arrays and objects
      if (type === "array" || type === "object") {
        try {
          value = JSON.parse(value)
        } catch (e) {
          console.error(`Failed to parse ${type} for ${name}:`, e)
          return
        }
      }

      // Convert numbers and booleans
      if (type === "integer" || type === "number") {
        value = Number(value)
      } else if (type === "boolean") {
        value = value === "true"
      }

      data[name] = value
    })

    return data
  }

  validateParameters() {
    const requiredParams = this.paramTargets.filter(
      (input) => input.dataset.paramRequired === "true"
    )

    for (const param of requiredParams) {
      if (!param.value.trim()) {
        param.classList.add("border-red-500")
        return false
      } else {
        param.classList.remove("border-red-500")
      }
    }

    return true
  }

  async makeRequest(data) {
    const csrfToken = document.querySelector('meta[name="csrf-token"]')?.content

    return fetch('/api/proxy', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': csrfToken
      },
      body: JSON.stringify({
        endpoint: this.urlValue,
        method: this.methodValue,
        params: data
      })
    })
  }

  async handleResponse(response, clientDuration) {
    // Show response container
    this.hideError()
    this.responseContainerTarget.classList.remove("hidden")

    const proxyResult = await response.json()

    // Use the actual API status code from the proxy envelope
    const actualStatus = proxyResult.status_code || response.status
    const statusOk = actualStatus >= 200 && actualStatus < 300

    // Update status badge
    const statusClass = statusOk ? "bg-green-100 text-green-800" : "bg-red-100 text-red-800"
    this.statusBadgeTarget.className = `inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium ${statusClass}`
    this.statusBadgeTarget.textContent = actualStatus.toString()

    // Update response time (prefer server-measured, fall back to client)
    const displayTime = proxyResult.response_time_ms || clientDuration
    this.responseTimeTarget.textContent = `${displayTime}ms`

    // Display the actual API response body (unwrap proxy envelope)
    const body = proxyResult.data ?? proxyResult.error ?? proxyResult
    this.responseBodyTarget.textContent = JSON.stringify(body, null, 2)
  }

  showLoading() {
    this.submitButtonTarget.disabled = true
    this.submitTextTarget.textContent = "Sending..."
    this.loadingTarget.classList.remove("hidden")
    this.hideError()
    this.responseContainerTarget.classList.add("hidden")
  }

  hideLoading() {
    this.submitButtonTarget.disabled = false
    this.submitTextTarget.textContent = "Send Request"
    this.loadingTarget.classList.add("hidden")
  }

  showError(message) {
    this.errorMessageTarget.textContent = message
    this.errorContainerTarget.classList.remove("hidden")
    this.responseContainerTarget.classList.add("hidden")
  }

  hideError() {
    this.errorContainerTarget.classList.add("hidden")
  }
}
