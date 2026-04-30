# frozen_string_literal: true

class ApplicationMailer < ActionMailer::Base
  default from: ENV.fetch("SMTP_FROM_EMAIL", "noreply@requiems.xyz")
  layout "mailer"
end
