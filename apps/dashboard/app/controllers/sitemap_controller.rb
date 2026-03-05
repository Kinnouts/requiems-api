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
end
