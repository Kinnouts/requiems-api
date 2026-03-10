# frozen_string_literal: true

class LocaleController < ApplicationController
  def update
    locale = params[:locale].presence
    locale = nil unless I18n.available_locales.map(&:to_s).include?(locale)

    if user_signed_in?
      current_user.update_column(:locale, locale)
      I18n.locale = locale&.to_sym || http_accept_language.compatible_language_from(I18n.available_locales) || I18n.default_locale
    else
      I18n.locale = locale&.to_sym || I18n.default_locale
    end

    redirect_back fallback_location: root_path
  end
end
