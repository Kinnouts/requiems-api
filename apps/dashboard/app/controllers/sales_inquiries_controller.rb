# frozen_string_literal: true

class SalesInquiriesController < ApplicationController
  def new
    # Render form
  end

  def create
    if valid_inquiry?
      # Send email to observers
      SalesMailer.enterprise_inquiry(inquiry_params).deliver_now

      redirect_to root_path, notice: "Thank you! Our sales team will contact you within 24 hours."
    else
      flash.now[:alert] = "Please fill in all required fields."
      render :new, status: :unprocessable_entity
    end
  rescue StandardError => e
    Rails.logger.error "Failed to send sales inquiry email: #{e.message}"
    flash.now[:alert] = "Sorry, there was an error submitting your inquiry. Please try again later."
    render :new, status: :unprocessable_entity
  end

  private

  def inquiry_params
    params.require(:inquiry).permit(:name, :email, :company, :message)
  end

  def valid_inquiry?
    params[:inquiry].present? &&
      params[:inquiry][:name].present? &&
      params[:inquiry][:email].present? &&
      params[:inquiry][:company].present? &&
      valid_email?(params[:inquiry][:email])
  end

  def valid_email?(email)
    email =~ URI::MailTo::EMAIL_REGEXP
  end
end
