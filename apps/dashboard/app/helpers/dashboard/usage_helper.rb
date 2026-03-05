# frozen_string_literal: true
module Dashboard::UsageHelper
  RANGE_DEFAULTS = { "7d" => "Last 7 Days", "30d" => "Last 30 Days" }.freeze

  def range_button_classes(key)
    active = params[:range] == key ||
             (params[:range].blank? && @range_label == RANGE_DEFAULTS[key])
    base = "px-4 py-2 rounded-lg text-sm font-medium transition-colors"
    state = active ? "bg-blue-600 text-white" : "bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600"
    "#{base} #{state}"
  end
end
