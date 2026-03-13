# frozen_string_literal: true

module ApisHelper
  # Load the API catalog from YAML config
  def api_catalog
    @api_catalog ||= YAML.load_file(Rails.root.join("config", "api_catalog.yml"))
  end

  # Get all categories
  def api_categories
    api_catalog["categories"]
  end

  # Get all APIs
  def all_apis
    api_catalog["apis"]
  end

  # Get APIs by category
  def apis_by_category(category_id)
    all_apis.select { |api| Array(api["categories"]).include?(category_id) }
  end

  # Get category by ID
  def find_category(category_id)
    api_categories.find { |cat| cat["id"] == category_id }
  end

  # Get API by ID
  def find_api(api_id)
    all_apis.find { |api| api["id"] == api_id }
  end

  # Get live APIs only
  def live_apis
    all_apis.select { |api| api["status"] == "live" }
  end

  # Get status badge variant
  def api_status_variant(status)
    case status
    when "live"
      "success"
    when "beta"
      "warning"
    when "deprecated"
      "danger"
    else
      "default"
    end
  end

  # Get status badge text
  def api_status_text(status)
    case status
    when "live"
      "Live"
    when "beta"
      "Beta"
    when "deprecated"
      "Deprecated"
    else
      status.to_s.titleize
    end
  end

  # Load API documentation from YAML
  def api_documentation(api_id)
    doc_path = Rails.root.join("config", "api_docs", "#{api_id}.yml")
    return nil unless File.exist?(doc_path)

    @api_docs ||= {}
    @api_docs[api_id] ||= YAML.load_file(doc_path)
  end

  # Get popular APIs (for "Most Popular" section)
  def popular_apis
    live_apis.select { |api| api["popular"] == true }
  end

  # Group live APIs by category (returns hash; an API appears under each of its categories)
  def apis_grouped_by_category
    live_apis.each_with_object({}) do |api, hash|
      Array(api["categories"]).each do |cat|
        (hash[cat] ||= []) << api
      end
    end
  end

  # Get categories that have live APIs
  def categories_with_apis
    api_categories.select do |category|
      !category["coming_soon"] && apis_by_category(category["id"]).any?
    end
  end

  def api_documentation_as_text(doc)
    lines = []
    lines << doc["api_name"]
    lines << "=" * doc["api_name"].length
    lines << ""
    lines << doc["description"]
    lines << ""
    lines << "Base URL: #{doc["base_url"]}"
    lines << ""

    doc["endpoints"]&.each do |ep|
      lines << "#{ep["method"]} #{ep["path"]}"
      lines << "-" * "#{ep["method"]} #{ep["path"]}".length
      lines << ep["description"] if ep["description"].present?
      lines << ""

      if ep["parameters"].present?
        lines << "Parameters:"
        ep["parameters"].each do |p|
          req = p["required"] ? " (required)" : " (optional)"
          lines << "  #{p["name"]} [#{p["type"]}]#{req} - #{p["description"]}"
        end
        lines << ""
      end

      if ep["response_example"].present?
        lines << "Response example:"
        lines << ep["response_example"].strip
        lines << ""
      end
    end

    lines.join("\n")
  end

  def api_documentation_as_markdown(doc)
    lines = []
    lines << "# #{doc["api_name"]}"
    lines << ""
    lines << doc["description"]
    lines << ""
    lines << "**Base URL:** `#{doc["base_url"]}`"
    lines << ""

    doc["endpoints"]&.each do |ep|
      lines << "## `#{ep["method"]} #{ep["path"]}`"
      lines << ""
      lines << ep["description"] if ep["description"].present?
      lines << ""

      if ep["parameters"].present?
        lines << "### Parameters"
        lines << ""
        lines << "| Name | Type | Required | Description |"
        lines << "|------|------|----------|-------------|"
        ep["parameters"].each do |p|
          req = p["required"] ? "Yes" : "No"
          lines << "| `#{p["name"]}` | #{p["type"]} | #{req} | #{p["description"]} |"
        end
        lines << ""
      end

      if ep["request_example"].present?
        lines << "### Request"
        lines << ""
        lines << "```json"
        lines << ep["request_example"].strip
        lines << "```"
        lines << ""
      end

      if ep["response_example"].present?
        lines << "### Response"
        lines << ""
        lines << "```json"
        lines << ep["response_example"].strip
        lines << "```"
        lines << ""
      end

      next unless ep["code_examples"].present?

      lines << "### Code Examples"
      lines << ""
      ep["code_examples"].each do |lang, code|
        lines << "**#{lang.capitalize}**"
        lines << ""
        lines << "```#{lang}"
        lines << code.strip
        lines << "```"
        lines << ""
      end
    end

    lines.join("\n")
  end

  CATEGORY_ICON_SVGS = {
    "validation"     => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75 11.25 15 15 9.75m-3-7.036A11.959 11.959 0 0 1 3.598 6 11.99 11.99 0 0 0 3 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285Z" /></svg>',
    "text"           => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z" /></svg>',
    "technology"     => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M17.25 6.75 22.5 12l-5.25 5.25m-10.5 0L1.5 12l5.25-5.25m7.5-3-4.5 16.5" /></svg>',
    "media"          => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z" /></svg>',
    "ai_vision"      => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M9.813 15.904 9 18.75l-.813-2.846a4.5 4.5 0 0 0-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 0 0 3.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 0 0 3.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 0 0-3.09 3.09ZM18.259 8.715 18 9.75l-.259-1.035a3.375 3.375 0 0 0-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 0 0 2.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 0 0 2.456 2.456L21.75 6l-1.035.259a3.375 3.375 0 0 0-2.456 2.456ZM16.894 20.567 16.5 21.75l-.394-1.183a2.25 2.25 0 0 0-1.423-1.423L13.5 18.75l1.183-.394a2.25 2.25 0 0 0 1.423-1.423l.394-1.183.394 1.183a2.25 2.25 0 0 0 1.423 1.423l1.183.394-1.183.394a2.25 2.25 0 0 0-1.423 1.423Z" /></svg>',
    "places"         => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" /><path stroke-linecap="round" stroke-linejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z" /></svg>',
    "entertainment"  => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="m15.75 10.5 4.72-4.72a.75.75 0 0 1 1.28.53v11.38a.75.75 0 0 1-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 0 0 2.25-2.25v-9a2.25 2.25 0 0 0-2.25-2.25h-9A2.25 2.25 0 0 0 2.25 7.5v9a2.25 2.25 0 0 0 2.25 2.25Z" /></svg>',
    "finance"        => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 0 0 2.25-2.25V6.75A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25v10.5A2.25 2.25 0 0 0 4.5 19.5Z" /></svg>',
    "security"       => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z" /></svg>',
    "health"         => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" /></svg>',
    "tax"            => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M9 14.25l6-6m4.5-3.493V21.75l-3.75-1.5-3.75 1.5-3.75-1.5-3.75 1.5V4.757c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0 1 11.186 0c1.1.128 1.907 1.077 1.907 2.185ZM9.75 9h.008v.008H9.75V9Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm4.125 4.5h.008v.008h-.008V13.5Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z" /></svg>',
    "transportation" => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M8.25 18.75a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m3 0h6m-9 0H3.375a1.125 1.125 0 0 1-1.125-1.125V14.25m17.25 4.5a1.5 1.5 0 0 1-3 0m3 0a1.5 1.5 0 0 0-3 0m3 0h1.125c.621 0 1.129-.504 1.09-1.124a17.902 17.902 0 0 0-3.213-9.193 2.056 2.056 0 0 0-1.58-.86H14.25M16.5 18.75h-2.25m0-11.177v-.958c0-.568-.422-1.048-.987-1.106a48.554 48.554 0 0 0-10.026 0 1.106 1.106 0 0 0-.987 1.106v7.635m12-6.677v6.677m0 4.5v-4.5m0 0h-12" /></svg>',
    "data"           => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" /></svg>',
    "nature"         => '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6"><circle cx="8.5" cy="6.5" r="1.5" fill="currentColor" stroke="none"/><circle cx="15.5" cy="6.5" r="1.5" fill="currentColor" stroke="none"/><circle cx="5.5" cy="11" r="1.25" fill="currentColor" stroke="none"/><circle cx="18.5" cy="11" r="1.25" fill="currentColor" stroke="none"/><path stroke-linecap="round" stroke-linejoin="round" d="M12 22c-3.5 0-6-2.5-6-5 0-1.5.75-2.75 1.75-3.5L9 12.5h6l1.25 1c1 .75 1.75 2 1.75 3.5 0 2.5-2.5 5-6 5Z"/></svg>'
  }.freeze

  CATEGORY_COLORS = {
    "validation"     => "bg-green-100 dark:bg-green-900 text-green-600 dark:text-green-300",
    "text"           => "bg-purple-100 dark:bg-purple-900 text-purple-600 dark:text-purple-300",
    "technology"     => "bg-cyan-100 dark:bg-cyan-900 text-cyan-600 dark:text-cyan-300",
    "media"          => "bg-pink-100 dark:bg-pink-900 text-pink-600 dark:text-pink-300",
    "ai_vision"      => "bg-violet-100 dark:bg-violet-900 text-violet-600 dark:text-violet-300",
    "places"         => "bg-emerald-100 dark:bg-emerald-900 text-emerald-600 dark:text-emerald-300",
    "entertainment"  => "bg-orange-100 dark:bg-orange-900 text-orange-500 dark:text-orange-300",
    "finance"        => "bg-amber-100 dark:bg-amber-900 text-amber-600 dark:text-amber-300",
    "security"       => "bg-teal-100 dark:bg-teal-900 text-teal-600 dark:text-teal-300",
    "health"         => "bg-rose-100 dark:bg-rose-900 text-rose-600 dark:text-rose-300",
    "tax"            => "bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300",
    "transportation" => "bg-sky-100 dark:bg-sky-900 text-sky-600 dark:text-sky-300",
    "data"           => "bg-indigo-100 dark:bg-indigo-900 text-indigo-600 dark:text-indigo-300",
    "nature"         => "bg-lime-100 dark:bg-lime-900 text-lime-600 dark:text-lime-300"
  }.freeze

  def category_icon_svg(category_id, size: "w-6 h-6")
    svg = CATEGORY_ICON_SVGS[category_id] || CATEGORY_ICON_SVGS["data"]
    svg.gsub('class="w-6 h-6"', "class=\"#{size}\"").html_safe
  end

  def category_color_classes(category_id, hover: false)
    base = CATEGORY_COLORS[category_id] || CATEGORY_COLORS["data"]
    return base unless hover

    hover_map = {
      "validation"     => "group-hover:bg-green-200 dark:group-hover:bg-green-800",
      "text"           => "group-hover:bg-purple-200 dark:group-hover:bg-purple-800",
      "technology"     => "group-hover:bg-cyan-200 dark:group-hover:bg-cyan-800",
      "media"          => "group-hover:bg-pink-200 dark:group-hover:bg-pink-800",
      "ai_vision"      => "group-hover:bg-violet-200 dark:group-hover:bg-violet-800",
      "places"         => "group-hover:bg-emerald-200 dark:group-hover:bg-emerald-800",
      "entertainment"  => "group-hover:bg-orange-200 dark:group-hover:bg-orange-800",
      "finance"        => "group-hover:bg-amber-200 dark:group-hover:bg-amber-800",
      "security"       => "group-hover:bg-teal-200 dark:group-hover:bg-teal-800",
      "health"         => "group-hover:bg-rose-200 dark:group-hover:bg-rose-800",
      "tax"            => "group-hover:bg-slate-200 dark:group-hover:bg-slate-700",
      "transportation" => "group-hover:bg-sky-200 dark:group-hover:bg-sky-800",
      "data"           => "group-hover:bg-indigo-200 dark:group-hover:bg-indigo-800",
      "nature"         => "group-hover:bg-lime-200 dark:group-hover:bg-lime-800"
    }
    "#{base} #{hover_map[category_id] || hover_map["data"]}"
  end

  def open_in_claude_url(documentation)
    doc_url = "#{request.base_url}/apis/#{documentation['api_id']}/index.md"
    prompt = "Read this page from the Requiems API docs: #{doc_url} and help me integrate this API into my project."
    "https://claude.ai/new?q=#{ERB::Util.url_encode(prompt)}"
  end

  def open_in_chatgpt_url(documentation)
    doc_url = "#{request.base_url}/apis/#{documentation['api_id']}/index.md"
    prompt = "Read this page from the Requiems API docs: #{doc_url} and help me integrate this API into my project."
    "https://chatgpt.com/?q=#{ERB::Util.url_encode(prompt)}"
  end
end
