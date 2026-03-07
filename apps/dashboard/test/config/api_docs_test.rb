# frozen_string_literal: true

require "test_helper"
require "yaml"

# Validates every YAML file under config/api_docs/ against the expected schema.
#
# Schema rules:
#   Top-level required: api_id, api_name, description, base_url, endpoints
#   endpoints[]: name, method (GET/POST/PUT/PATCH/DELETE), path (/v1/…), description
#   parameters[]: name, type, required (bool), location (path|query|body), description
#     - "location" is the required key. Using "in" (OpenAPI) is not allowed.
class ApiDocsSchemaTest < ActiveSupport::TestCase
  DOCS_DIR = Rails.root.join("config/api_docs")
  VALID_METHODS = %w[GET POST PUT PATCH DELETE].freeze
  VALID_LOCATIONS = %w[path query body].freeze

  Dir[DOCS_DIR.join("*.yml")].sort.each do |file|
    basename = File.basename(file)
    doc = YAML.load_file(file)

    test "#{basename}: has required top-level fields" do
      %w[api_id api_name description base_url endpoints].each do |field|
        assert doc.key?(field), "#{basename}: missing top-level field '#{field}'"
        assert doc[field].present?, "#{basename}: '#{field}' must not be blank"
      end
    end

    test "#{basename}: endpoints is a non-empty array" do
      assert_kind_of Array, doc["endpoints"], "#{basename}: 'endpoints' must be an array"
      assert doc["endpoints"].any?, "#{basename}: 'endpoints' must have at least one entry"
    end

    doc["endpoints"]&.each_with_index do |ep, i|
      label = "#{basename} endpoint[#{i}] (#{ep["name"] || "unnamed"})"

      test "#{label}: has required fields" do
        %w[name method path description].each do |field|
          assert ep.key?(field), "#{label}: missing field '#{field}'"
          assert ep[field].present?, "#{label}: '#{field}' must not be blank"
        end
      end

      test "#{label}: method is valid" do
        assert_includes VALID_METHODS, ep["method"],
          "#{label}: method '#{ep["method"]}' must be one of #{VALID_METHODS.join(", ")}"
      end

      test "#{label}: path starts with /v1/" do
        assert ep["path"].to_s.start_with?("/v1/"),
          "#{label}: path '#{ep["path"]}' must start with /v1/"
      end

      (ep["parameters"] || []).each_with_index do |param, pi|
        plabel = "#{label} param[#{pi}] (#{param["name"] || "unnamed"})"

        test "#{plabel}: has required fields" do
          %w[name type description].each do |field|
            assert param.key?(field), "#{plabel}: missing field '#{field}'"
            assert param[field].present?, "#{plabel}: '#{field}' must not be blank"
          end
          assert param.key?("required"), "#{plabel}: missing field 'required'"
          assert_includes [ true, false ], param["required"],
            "#{plabel}: 'required' must be a boolean"
        end

        test "#{plabel}: uses 'location' not 'in'" do
          assert param.key?("location"),
            "#{plabel}: must use 'location:' (not 'in:'). " \
            "Valid values: #{VALID_LOCATIONS.join(", ")}"
          refute param.key?("in"),
            "#{plabel}: found 'in:' key — use 'location:' instead"
        end

        test "#{plabel}: location is valid" do
          assert_includes VALID_LOCATIONS, param["location"],
            "#{plabel}: location '#{param["location"]}' must be one of #{VALID_LOCATIONS.join(", ")}"
        end
      end
    end
  end
end
