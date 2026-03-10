import { Controller } from "@hotwired/stimulus";

export default class extends Controller {
  static targets = ["option"];

  select(event) {
    const selected = event.currentTarget;
    this.optionTargets.forEach((option) => {
      const isSelected = option === selected;
      option.classList.toggle("border-blue-500", isSelected);
      option.classList.toggle("bg-blue-50", isSelected);
      option.classList.toggle("dark:bg-blue-900/30", isSelected);
      option.classList.toggle("text-blue-700", isSelected);
      option.classList.toggle("dark:text-blue-300", isSelected);
      option.classList.toggle("border-gray-300", !isSelected);
      option.classList.toggle("dark:border-gray-600", !isSelected);
      option.classList.toggle("text-gray-700", !isSelected);
      option.classList.toggle("dark:text-gray-300", !isSelected);
    });
  }
}
