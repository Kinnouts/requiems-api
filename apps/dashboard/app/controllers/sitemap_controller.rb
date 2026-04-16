# frozen_string_literal: true

class SitemapController < ApplicationController
  include ApisHelper

  before_action :set_response_content_type

  def sitemap
    expires_in 5.minutes, public: true
    @apis = live_apis
    @last_modified = Time.current.beginning_of_day
    respond_to { |f| f.xml }
  end

  def llms
    expires_in 5.minutes, public: true
    @apis = live_apis
    respond_to { |f| f.text }
  end

  def llms_full
    expires_in 5.minutes, public: true
    @docs = live_apis.filter_map { |api| api_documentation(api["id"]) }
    content = render_to_string(template: "sitemap/llms-full", formats: [ :text ], layout: false)
    send_data content, filename: "llms-full.txt", type: "text/plain", disposition: "attachment"
  end

  def api_doc
    expires_in 5.minutes, public: true
    @api = find_api(params[:id])
    return head :not_found unless @api

    @doc = api_documentation(params[:id])
    return head :not_found unless @doc

    respond_to { |f| f.text }
  end

  private

  def set_response_content_type
    case action_name
    when "sitemap"
      response.content_type = "application/xml"
    when "llms", "llms_full", "api_doc"
      response.content_type = "text/plain"
    end
  end
end
