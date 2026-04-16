# frozen_string_literal: true

class SitemapController < ApplicationController
  include ApisHelper

  def sitemap
    expires_in 5.minutes, public: true
    @apis = live_apis
    @last_modified = Time.current.beginning_of_day
    render "sitemap/sitemap", formats: [ :xml ], layout: false, content_type: "application/xml"
  end

  def llms
    expires_in 5.minutes, public: true
    @apis = live_apis
    render "sitemap/llms", formats: [ :text ], layout: false, content_type: "text/plain"
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

    render "sitemap/api_doc", formats: [ :text ], layout: false, content_type: "text/plain"
  end
end
