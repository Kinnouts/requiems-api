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
end
