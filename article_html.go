package main

import (
	"fmt"
)

func ReadableArticleHTML() string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<meta content="text/html;charset=utf-8" http-equiv="Content-Type">
				<meta content="utf-8" http-equiv="encoding">
				<title>{{.Title}}</title>
				<style>%s</style>
			</head>
			<body>
				<header>
					<h1>
						<a href="{{.Address}}">{{.Title}}</a>
						<div class="archived-at">Archived at {{.ArchivedAt}}</div>
					</h1>
					<img src="{{.Image}}">
					<figcaption>{{.Excerpt}}</figcaption>
					<small>{{.Byline}} • {{.SiteName}} • {{.ReadingTime}} minutes</small>
				</header>
				<article>{{.Content}}</article>
			</body>
		</html>
		`, Styles())
}

func NonReadableArticleHTML() string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
	<head>
		<meta content="text/html;charset=utf-8" http-equiv="Content-Type">
		<meta content="utf-8" http-equiv="encoding">
		<title>{{.Address}}</title>
		<style>%s</style>
	</head>
	<body>
	  <header>
			<h1><a href="{{.Address}}">{{.Address}}</a></h1>
			<figcaption>{{.Error}}</figcaption>
	  </header>
	</body>
</html>
	`, Styles())
}
