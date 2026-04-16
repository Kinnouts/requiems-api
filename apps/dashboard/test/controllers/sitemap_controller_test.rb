# frozen_string_literal: true

require "test_helper"

class SitemapControllerTest < ActionDispatch::IntegrationTest
  test "sitemap is publicly cacheable" do
    get "/sitemap.xml"

    assert_response :success
    assert_match %r{\Aapplication/xml(?:;.*)?\z}, response.headers["Content-Type"]
    assert_includes response.headers["Cache-Control"], "public"
    assert_includes response.headers["Cache-Control"], "max-age=300"
  end

  test "llms file is publicly cacheable" do
    get "/llms.txt"

    assert_response :success
    assert_match %r{\Atext/plain(?:;.*)?\z}, response.headers["Content-Type"]
    assert_includes response.headers["Cache-Control"], "public"
    assert_includes response.headers["Cache-Control"], "max-age=300"
  end

  test "llms full file is text/plain" do
    get "/llms-full.txt"

    assert_response :success
    assert_match %r{\Atext/plain(?:;.*)?\z}, response.headers["Content-Type"]
  end
end
