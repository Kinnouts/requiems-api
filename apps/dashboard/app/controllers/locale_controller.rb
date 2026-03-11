# frozen_string_literal: true

class LocaleController < ApplicationController
  def update
    locale = params[:locale].presence
    locale = nil unless I18n.available_locales.map(&:to_s).include?(locale)

    current_user.update(locale: locale) if user_signed_in?

    I18n.locale = locale&.to_sym || I18n.default_locale

    redirect_to redirect_path
  end

  private

  def redirect_path
    ref = request.referer
    return root_path unless ref

    uri = URI.parse(ref)
    # Strip existing locale prefix from path
    locale_pattern = I18n.available_locales.map { |l| Regexp.escape(l.to_s) }.join("|")
    path = uri.path.sub(%r{\A/(#{locale_pattern})(?=/|\z)}, "")
    path = "/" if path.blank?

    # Prepend new locale prefix if non-default
    path = "/#{I18n.locale}#{path}" if I18n.locale != I18n.default_locale

    # Preserve query string and fragment
    path = "#{path}?#{uri.query}" if uri.query.present?
    path = "#{path}##{uri.fragment}" if uri.fragment.present?

    path
  rescue URI::InvalidURIError
    root_path
  end
end
