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
    all_apis.select { |api| api["category"] == category_id }
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

  # Group live APIs by category (returns hash)
  def apis_grouped_by_category
    live_apis.group_by { |api| api["category"] }
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
end
