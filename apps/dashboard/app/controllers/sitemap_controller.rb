# frozen_string_literal: true

class SitemapController < ApplicationController
  include ApisHelper

  def sitemap
    @apis = live_apis
    @last_modified = Time.current.beginning_of_day
    respond_to { |f| f.xml }
  end

  def llms
    @apis = live_apis
    respond_to { |f| f.text }
  end

  def llms_full
    @docs = live_apis.filter_map { |api| api_documentation(api["id"]) }
    content = render_to_string(template: "sitemap/llms-full", formats: [ :text ], layout: false)
    send_data content, filename: "llms-full.txt", type: "text/plain", disposition: "attachment"
  end

  def api_doc
    @api = find_api(params[:id])
    return head :not_found unless @api

    @doc = api_documentation(params[:id])
    return head :not_found unless @doc

    respond_to { |f| f.text }
  end
end
