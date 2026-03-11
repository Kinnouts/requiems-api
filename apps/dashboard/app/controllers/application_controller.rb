# frozen_string_literal: true

class ApplicationController < ActionController::Base
  include Pagy::Backend
  include HttpAcceptLanguage::EasyAccess

  # Only allow modern browsers supporting webp images, web push, badges, import maps, CSS nesting, and CSS :has.
  allow_browser versions: :modern

  # Changes to the importmap will invalidate the etag for HTML responses
  stale_when_importmap_changes

  before_action :set_locale

  private

  def set_locale
    requested = params[:locale].presence
    validated = requested if I18n.available_locales.map(&:to_s).include?(requested)
    I18n.locale = validated ||
                  current_user_locale ||
                  http_accept_language.compatible_language_from(I18n.available_locales) ||
                  I18n.default_locale

    # Redirect to the localized URL when locale was resolved from user preference or
    # Accept-Language (not from the URL itself) so that the URL always reflects the
    # page language — keeps canonical tags and hreflang accurate.
    return if validated.present? || request.path.start_with?("/dashboard", "/admin", "/users")

    if I18n.locale != I18n.default_locale && params[:locale].blank?
      redirect_to url_for(locale: I18n.locale), status: :moved_permanently, allow_other_host: false
    end
  end

  def current_user_locale
    current_user&.locale&.presence&.to_sym
  end

  def default_url_options
    { locale: I18n.locale == I18n.default_locale ? nil : I18n.locale }
  end
end
