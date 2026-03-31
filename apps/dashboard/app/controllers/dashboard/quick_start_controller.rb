# frozen_string_literal: true

class Dashboard::QuickStartController < ApplicationController
  before_action :authenticate_user!
  layout "dashboard"

  def index
  end
end
