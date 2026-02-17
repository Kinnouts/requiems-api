class HomeController < ApplicationController
  def index
    @categories = YAML.load_file(Rails.root.join("config", "api_catalog.yml"))["categories"]
  end

  def docs
  end

  def pricing
  end

  def blog
  end

  def status
  end

  def glossary
  end

  def error_codes
  end

  def faq
  end

  def team
  end
end
