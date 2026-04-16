# frozen_string_literal: true

# Pretty-print the sitemap after sitemap_generator writes it.
# sitemap_generator minifies all XML onto one line; this makes Chrome render it
# as readable XML instead of a wall of text.
Rake::Task["sitemap:refresh"].enhance do
  require "rexml/document"

  path = Rails.root.join("public", "sitemap.xml")
  next unless path.exist?

  doc = REXML::Document.new(path.read)
  fmt = REXML::Formatters::Pretty.new(2)
  fmt.compact = true
  out = +""
  fmt.write(doc, out)
  path.write("#{out}\n")

  puts "sitemap: pretty-printed #{path}"
end
