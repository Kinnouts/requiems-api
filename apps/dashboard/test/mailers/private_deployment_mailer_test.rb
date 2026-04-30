# frozen_string_literal: true

require "test_helper"

class PrivateDeploymentMailerTest < ActionMailer::TestCase
  def setup
    @user = create_user(email: "deployments@example.com")
    @pdr = PrivateDeploymentRequest.create!(
      user: @user,
      company: "Acme Corp",
      contact_name: "Jordan",
      contact_email: @user.email,
      server_tier: "starter",
      billing_cycle: "monthly",
      status: "active",
      selected_services: %w[email],
      monthly_price_cents: 20_000,
      subdomain_slug: "acme-private",
      tenant_secret: "s" * 32
    )
  end

  test "request_received sends customer confirmation" do
    email = PrivateDeploymentMailer.request_received(@pdr)

    assert_equal [ @pdr.contact_email ], email.to
    assert_equal "Your Private Deployment Request — We're On It", email.subject
    assert_includes email.body.encoded, "https://requiems.xyz/docs"
  end

  test "admin_notification sends admin alert with reply_to" do
    email = PrivateDeploymentMailer.admin_notification(@pdr)

    assert_equal OBSERVER_EMAILS, email.to
    assert_equal [ @pdr.contact_email ], email.reply_to
    assert_equal "New Private Deployment Request: Acme Corp", email.subject
    assert_includes email.body.encoded, "/admin/private_deployments/#{@pdr.id}"
  end

  test "deployment_ready includes live url in subject and body" do
    email = PrivateDeploymentMailer.deployment_ready(@pdr)

    assert_equal [ @pdr.contact_email ], email.to
    assert_equal "Your Private API is Live — https://acme-private.requiems.xyz", email.subject
    assert_includes email.body.encoded, @pdr.live_url
  end
end
