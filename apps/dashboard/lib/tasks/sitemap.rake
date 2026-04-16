# frozen_string_literal: true

# Pretty-print and add XSLT stylesheet reference after sitemap_generator writes.
# sitemap_generator minifies onto one line; REXML formats it, then we inject the
# <?xml-stylesheet?> PI so browsers render a styled HTML view instead of raw XML.
Rake::Task["sitemap:refresh"].enhance do
  require "rexml/document"

  path = Rails.root.join("public", "sitemap.xml")
  next unless path.exist?

  doc = REXML::Document.new(path.read)
  fmt = REXML::Formatters::Pretty.new(2)
  fmt.compact = true
  out = +""
  fmt.write(doc, out)

  # Inject XSL stylesheet PI immediately after the XML declaration so browsers
  # render the sitemap as a styled HTML page rather than raw XML source.
  out.sub!(
    "<?xml version='1.0' encoding='UTF-8'?>",
    "<?xml version='1.0' encoding='UTF-8'?>\n<?xml-stylesheet type='text/xsl' href='/sitemap.xsl'?>"
  )

  path.write("#{out}\n")
  puts "sitemap: generated #{path}"
end
