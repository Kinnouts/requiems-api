# frozen_string_literal: true

class ApisController < ApplicationController
  include ApisHelper

  def index
    @categories = api_categories
    @apis = all_apis
    @selected_category = params[:category]

    # Filter by category if specified
    if @selected_category.present?
      @apis = apis_by_category(@selected_category)
      @category = find_category(@selected_category)
    end
  end

  def show
    @api = find_api(params[:id])

    if @api.nil?
      redirect_to apis_path, alert: "API not found"
      return
    end

    @category = find_category(@api["category"])
    @categories = api_categories
  end
end
