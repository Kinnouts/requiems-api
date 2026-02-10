import { Controller } from "@hotwired/stimulus"

// Interactive API playground for testing endpoints
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
    "responseHeaders",
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
      const duration = Math.round(endTime - this.startTime)

      await this.handleResponse(response, duration)
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
    const options = {
      method: this.methodValue,
      headers: {
        "Content-Type": "application/json"
      }
    }

    // Add request body for POST/PUT/PATCH
    if (["POST", "PUT", "PATCH"].includes(this.methodValue)) {
      options.body = JSON.stringify(data)
    }

    // For GET requests with query parameters, append to URL
    if (this.methodValue === "GET" && Object.keys(data).length > 0) {
      const params = new URLSearchParams(data)
      const url = `${this.urlValue}?${params.toString()}`
      return fetch(url, options)
    }

    return fetch(this.urlValue, options)
  }

  async handleResponse(response, duration) {
    // Show response container
    this.hideError()
    this.responseContainerTarget.classList.remove("hidden")

    // Update status badge
    const statusClass = response.ok ? "bg-green-100 text-green-800" : "bg-red-100 text-red-800"
    this.statusBadgeTarget.className = `inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium ${statusClass}`
    this.statusBadgeTarget.textContent = `${response.status} ${response.statusText}`

    // Update response time
    this.responseTimeTarget.textContent = `${duration}ms`

    // Update response headers
    const headers = {}
    response.headers.forEach((value, key) => {
      headers[key] = value
    })
    this.responseHeadersTarget.textContent = JSON.stringify(headers, null, 2)

    // Update response body
    try {
      const contentType = response.headers.get("content-type")
      let body

      if (contentType && contentType.includes("application/json")) {
        body = await response.json()
        this.responseBodyTarget.textContent = JSON.stringify(body, null, 2)
      } else {
        body = await response.text()
        this.responseBodyTarget.textContent = body
      }
    } catch (error) {
      this.responseBodyTarget.textContent = "Error parsing response"
      console.error("Response parsing error:", error)
    }
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
