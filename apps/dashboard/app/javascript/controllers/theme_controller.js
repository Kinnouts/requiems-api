import { Controller } from "@hotwired/stimulus";

export default class extends Controller {
  static targets = ["sun", "moon"];

  connect() {
    this.boundUpdate = this.updateIcons.bind(this);
    window.addEventListener("theme:changed", this.boundUpdate);
    this.updateIcons();
  }

  disconnect() {
    window.removeEventListener("theme:changed", this.boundUpdate);
  }

  toggle() {
    const isDark = document.documentElement.classList.toggle("dark");
    localStorage.setItem("theme", isDark ? "dark" : "light");
    window.dispatchEvent(
      new CustomEvent("theme:changed", { detail: { isDark } }),
    );
    this.updateIcons();
  }

  updateIcons() {
    const isDark = document.documentElement.classList.contains("dark");
    if (this.hasSunTarget) this.sunTarget.classList.toggle("hidden", !isDark);
    if (this.hasMoonTarget) this.moonTarget.classList.toggle("hidden", isDark);
  }
}
