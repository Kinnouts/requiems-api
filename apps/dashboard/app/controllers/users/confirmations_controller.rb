# frozen_string_literal: true

class Users::ConfirmationsController < Devise::ConfirmationsController
  def new
    redirect_to root_path
  end

  def show
    return redirect_to root_path if params[:confirmation_token].blank?

    super
  end
end
