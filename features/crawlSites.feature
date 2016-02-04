Feature: Crawl Websites
	As a Site Admin
	I want a map of assets and inter page links for my website
	So that I can see the dependencies and relationships between my website pages

	Scenario: No arguments
		When I run the crawler with the following arguments ""
		Then I should see an error informing me "please provide a url as the 1st argument"
		And the exit code should be: 1

	Scenario: Invalid URL
		Given an invalid url "%^&"
		When I run the crawler with the following arguments "%^&"
		Then I should see an error informing me "unable to parse url"
		And the exit code should be: 1

	Scenario: Unfetchable URL
		Given a unfetchable url "http://localhost:40123"
		When I run the crawler with the following arguments "http://localhost:40123"
		Then I should see an error informing me "unable to fetch url"
		And the exit code should be: 1

	Scenario: Single page site with external links
		Given a webpage "index.html" containing:
			"""
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Index</title>
	</head>
	<body>
		<a href="http://google.com">Google</a>
		<a href="http://facebook.com">Facebook</a>
	</body>
</html>
			"""
		And the webpages are being hosted locally
		When I run the crawler with the locally hosted url
		Then I should see the following:
			"""
/
	Static Assets:
	Internal Links:
			"""
		And the exit code should be: 0

	Scenario: Single page site with static assets
		Given a webpage "index.html" containing:
			"""
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Index</title>
		<link rel="stylesheet" href="/static/css/style.css"/>
		<script src="/static/js/script.js"></script>
	</head>
	<body>
	</body>
</html>
			"""
		And the webpages are being hosted locally
		When I run the crawler with the locally hosted url
		Then I should see the following:
			"""
/
	Static Assets:
		CSS:
			http://127.0.0.1:\d+/static/css/style.css
		JS:
			http://127.0.0.1:\d+/static/js/script.js
	Internal Links:
			"""
		And the exit code should be: 0

	Scenario: Multi page site with internal links
		Given a webpage "index.html" containing:
			"""
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Index</title>
	</head>
	<body>
		<a href="/blog.html">Blog</a>
		<a href="/about">about</a>
	</body>
</html>
			"""
		And a webpage "blog.html" containing:
			"""
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Blog</title>
	</head>
	<body>
		<a href="/index.html">Home</a>
		<a href="/about">about</a>
	</body>
</html>
			"""
		And a webpage "about" containing:
			"""
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>About</title>
	</head>
	<body>
		<a href="/index.html">Home</a>
		<a href="/blog.html">Blog</a>
	</body>
</html>
			"""
		And the webpages are being hosted locally
		When I run the crawler with the locally hosted url
		Then I should see the following:
			"""
/
	Static Assets:
	Internal Links:
		http://127.0.0.1:\d+/(blog\.html|about)
		http://127.0.0.1:\d+/(blog\.html|about)
			"""
		And I should see the following:
			"""
/blog.html
	Static Assets:
	Internal Links:
		http://127.0.0.1:\d+/(index\.html|about)
		http://127.0.0.1:\d+/(index\.html|about)
			"""
		And I should see the following:
			"""
/index.html
	Static Assets:
	Internal Links:
		http://127.0.0.1:\d+/blog.html
		http://127.0.0.1:\d+/about
			"""
		And I should see the following:
			"""
/about
	Static Assets:
	Internal Links:
		http://127.0.0.1:\d+/(index\.html|blog\.html)
		http://127.0.0.1:\d+/(index\.html|blog\.html)
			"""

	Scenario: Multi page site with static assets, internal links and external links
		Given a webpage "index.html" containing:
			"""
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Index</title>
		<link rel="stylesheet" href="/static/css/style.css"/>
		<script src="/static/js/script.js"></script>
	</head>
	<body>
		<a href="/blog.html">Blog</a>
		<a href="/about">about</a>
	</body>
</html>
			"""
		And a webpage "blog.html" containing:
			"""
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Blog</title>
		<link rel="stylesheet" href="/static/css/style.css"/>
		<script src="/static/js/script.js"></script>
	</head>
	<body>
		<a href="/index.html">Home</a>
		<a href="/about">about</a>
	</body>
</html>
			"""
		And a webpage "about" containing:
			"""
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>About</title>
		<link rel="stylesheet" href="/static/css/style.css"/>
		<script src="/static/js/script.js"></script>
	</head>
	<body>
		<a href="/index.html">Home</a>
		<a href="/blog.html">Blog</a>
	</body>
</html>
			"""
		And the webpages are being hosted locally
		When I run the crawler with the locally hosted url
		Then I should see the following:
			"""
/
	Static Assets:
		CSS:
			http://127.0.0.1:\d+/static/css/style.css
		JS:
			http://127.0.0.1:\d+/static/js/script.js
	Internal Links:
		http://127.0.0.1:\d+/(blog\.html|about)
		http://127.0.0.1:\d+/(blog\.html|about)
			"""
		And I should see the following:
			"""
/blog.html
	Static Assets:
		CSS:
			http://127.0.0.1:\d+/static/css/style.css
		JS:
			http://127.0.0.1:\d+/static/js/script.js
	Internal Links:
		http://127.0.0.1:\d+/(index\.html|about)
		http://127.0.0.1:\d+/(index\.html|about)
			"""
		And I should see the following:
			"""
/about
	Static Assets:
		CSS:
			http://127.0.0.1:\d+/static/css/style.css
		JS:
			http://127.0.0.1:\d+/static/js/script.js
	Internal Links:
		http://127.0.0.1:\d+/(index\.html|blog\.html)
		http://127.0.0.1:\d+/(index\.html|blog\.html)
			"""
