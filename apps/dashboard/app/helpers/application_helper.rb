# frozen_string_literal: true

module ApplicationHelper
  include ApisHelper

  # Global search data for navbar search
  # Returns all searchable content (APIs, Examples, Pages) as a hash
  def global_search_data
    {
      apis: searchable_apis,
      examples: searchable_examples,
      pages: searchable_pages
    }
  end

  private

  def searchable_apis
    live_apis.map do |api|
      category = find_category(api["category"])
      {
        id: api["id"],
        title: api["name"],
        description: api["description"],
        url: api["documentation_url"],
        category: category["name"],
        category_icon: category["icon"],
        type: "api",
        tags: api["tags"] || [],
        endpoints_count: api["endpoints_count"]
      }
    end
  end

  def searchable_examples
    examples_data = YAML.load_file(Rails.root.join("config", "examples.yml"))
    examples_data["examples"].map do |example|
      {
        id: example["id"],
        title: example["title"],
        description: example["description"],
        url: "/examples/#{example['id']}",
        category: example["category"],
        type: "example",
        difficulty: example["difficulty"],
        technologies: example["technologies"]
      }
    end
  end

  def searchable_pages
    [
      {
        title: "Documentation",
        description: "Complete API documentation and integration guides",
        url: "/docs",
        type: "page",
        icon: "📚",
        tags: [ "help", "guide", "reference", "docs" ]
      },
      {
        title: "Pricing",
        description: "View pricing plans and choose the right plan for your needs",
        url: "/pricing",
        type: "page",
        icon: "💰",
        tags: [ "plans", "billing", "subscription", "cost" ]
      },
      {
        title: "API Reference",
        description: "Complete API reference with endpoints, parameters, and examples",
        url: "/api_reference",
        type: "page",
        icon: "📖",
        tags: [ "reference", "endpoints", "documentation" ]
      },
      {
        title: "Examples",
        description: "Browse code examples and tutorials for all APIs",
        url: "/examples",
        type: "page",
        icon: "💻",
        tags: [ "tutorials", "code", "samples", "examples" ]
      },
      {
        title: "FAQ",
        description: "Frequently asked questions and answers",
        url: "/faq",
        type: "page",
        icon: "❓",
        tags: [ "help", "questions", "support" ]
      },
      {
        title: "About",
        description: "Learn more about Requiems API and our mission",
        url: "/about",
        type: "page",
        icon: "ℹ️",
        tags: [ "company", "team", "mission" ]
      },
      {
        title: "Contact",
        description: "Get in touch with our team",
        url: "/contact",
        type: "page",
        icon: "📧",
        tags: [ "support", "help", "email", "message" ]
      },
      {
        title: "Changelog",
        description: "See what's new and what's changed in our APIs",
        url: "/changelog",
        type: "page",
        icon: "📝",
        tags: [ "updates", "releases", "versions", "news" ]
      },
      {
        title: "Status",
        description: "Check the status of our APIs and services",
        url: "/status",
        type: "page",
        icon: "🟢",
        tags: [ "uptime", "availability", "monitoring", "health" ]
      },
      {
        title: "Blog",
        description: "Read our latest blog posts and articles",
        url: "/blog",
        type: "page",
        icon: "📰",
        tags: [ "articles", "news", "updates", "posts" ]
      },
      {
        title: "Error Codes",
        description: "Reference for all API error codes and how to fix them",
        url: "/error_codes",
        type: "page",
        icon: "⚠️",
        tags: [ "errors", "troubleshooting", "debugging", "codes" ]
      },
      {
        title: "Glossary",
        description: "API terminology and definitions",
        url: "/glossary",
        type: "page",
        icon: "📓",
        tags: [ "terms", "definitions", "vocabulary", "reference" ]
      },
      {
        title: "Privacy Policy",
        description: "Read our privacy policy and data practices",
        url: "/privacy",
        type: "page",
        icon: "🔒",
        tags: [ "privacy", "data", "security", "policy" ]
      },
      {
        title: "Terms of Service",
        description: "Read our terms of service and usage policies",
        url: "/terms",
        type: "page",
        icon: "📄",
        tags: [ "terms", "legal", "policy", "agreement" ]
      }
    ]
  end
end
