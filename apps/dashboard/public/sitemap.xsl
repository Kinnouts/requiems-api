<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
  xmlns:sm="http://www.sitemaps.org/schemas/sitemap/0.9"
  xmlns:xhtml="http://www.w3.org/1999/xhtml"
  exclude-result-prefixes="sm xhtml">

  <xsl:output method="html" version="1.0" encoding="UTF-8" indent="yes"/>

  <xsl:template match="/">
    <html lang="en">
      <head>
        <meta charset="UTF-8"/>
        <meta name="viewport" content="width=device-width, initial-scale=1"/>
        <title>Sitemap &#8212; requiems.xyz</title>
        <style>
          *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
          body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
            color: #111;
            background: #fafafa;
            padding: 40px 20px;
          }
          .container { max-width: 1100px; margin: 0 auto; }
          header { margin-bottom: 32px; }
          header h1 { font-size: 1.4rem; font-weight: 600; margin-bottom: 6px; }
          header p { color: #555; font-size: 0.9rem; }
          header a { color: #2563eb; text-decoration: none; }
          header a:hover { text-decoration: underline; }
          table { width: 100%; border-collapse: collapse; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 4px rgba(0,0,0,0.08); }
          thead { background: #f0f0f0; }
          th {
            text-align: left;
            padding: 10px 16px;
            font-size: 0.75rem;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.06em;
            color: #444;
            border-bottom: 1px solid #ddd;
          }
          td {
            padding: 9px 16px;
            font-size: 0.85rem;
            border-bottom: 1px solid #f0f0f0;
            vertical-align: middle;
          }
          tbody tr:last-child td { border-bottom: none; }
          tbody tr:hover td { background: #f7f7f7; }
          td a { color: #2563eb; text-decoration: none; word-break: break-all; }
          td a:hover { text-decoration: underline; }
          .priority { font-variant-numeric: tabular-nums; }
          .freq { color: #555; }
          .lastmod { color: #555; white-space: nowrap; }
        </style>
      </head>
      <body>
        <div class="container">
          <header>
            <h1>XML Sitemap</h1>
            <p>
              <xsl:value-of select="count(sm:urlset/sm:url)"/>&#160;URLs &#8212;
              <a href="https://requiems.xyz">requiems.xyz</a>
            </p>
          </header>
          <table>
            <thead>
              <tr>
                <th>URL</th>
                <th>Priority</th>
                <th>Change&#160;Freq</th>
                <th>Last&#160;Modified</th>
              </tr>
            </thead>
            <tbody>
              <xsl:for-each select="sm:urlset/sm:url">
                <xsl:sort select="sm:priority" data-type="number" order="descending"/>
                <tr>
                  <td><a href="{sm:loc}"><xsl:value-of select="sm:loc"/></a></td>
                  <td class="priority"><xsl:value-of select="sm:priority"/></td>
                  <td class="freq"><xsl:value-of select="sm:changefreq"/></td>
                  <td class="lastmod"><xsl:value-of select="substring(sm:lastmod,1,10)"/></td>
                </tr>
              </xsl:for-each>
            </tbody>
          </table>
        </div>
      </body>
    </html>
  </xsl:template>
</xsl:stylesheet>
