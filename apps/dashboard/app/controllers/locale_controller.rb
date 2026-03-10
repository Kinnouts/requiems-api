# frozen_string_literal: true

class LocaleController < ApplicationController
  def update
    locale = params[:locale].presence
    locale = nil unless I18n.available_locales.map(&:to_s).include?(locale)

    current_user.update_column(:locale, locale) if user_signed_in?

    I18n.locale = locale&.to_sym || I18n.default_locale

    redirect_to redirect_path
  end

  private

  def redirect_path
    ref = request.referer
    return root_path unless ref

    uri = URI.parse(ref)
    # Strip existing locale prefix from path
    path = uri.path.sub(%r{\A/(en|es)(?=/|\z)}, "")
    path = "/" if path.blank?

    # Prepend new locale prefix if non-default
    path = "/#{I18n.locale}#{path}" if I18n.locale != I18n.default_locale

    path
  rescue URI::InvalidURIError
    root_path
  end
end
