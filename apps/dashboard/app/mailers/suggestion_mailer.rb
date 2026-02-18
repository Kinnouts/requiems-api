# frozen_string_literal: true

class SuggestionMailer < ApplicationMailer
  # Send new API suggestion notification to observers
  #
  # @param suggestion [Hash] The suggestion details
  #   @option suggestion [String] :api_name Name of the suggested API
  #   @option suggestion [String] :description Description of what the API should do
  #   @option suggestion [String] :use_case How the user would use the API
  #   @option suggestion [String] :email Optional user email for follow-up
  def new_suggestion(suggestion)
    @suggestion = suggestion
    @timestamp = Time.current

    mail(
      to: OBSERVER_EMAILS,
      subject: "New API Suggestion: #{suggestion[:api_name]}"
    )
  end
end
