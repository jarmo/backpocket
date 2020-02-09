package template

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
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<title>{{.Title}}</title>
				<style>%s</style>
				<style>%s</style>
			</head>
			<body>
				<header>
					<h1>{{.Title}}</h1>
					<div class="source-info">
						<div>Archived at {{.ArchivedAt}}</div>
						<div><a href="{{.Address}}">View Original</a></div>
					</div>
					<figure>
						<img src="{{.Image}}">
						<figcaption>{{.Excerpt}}</figcaption>
					</figure>
					<small>{{.Byline}} • {{.SiteName}} • {{.ReadingTime}} minutes</small>
				</header>
				<article>{{.Content}}</article>
			</body>
		</html>
		`, ModernNormalizeStyles(), Styles())
}

func NonReadableArticleHTML() string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
	<head>
		<meta content="text/html;charset=utf-8" http-equiv="Content-Type">
		<meta content="utf-8" http-equiv="encoding">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>{{.Address}}</title>
		<style>%s</style>
		<style>%s</style>
	</head>
	<body>
	  <header>
			<h1>Failed to archive</h1>
			<div class="source-info">
				<div>Tried at {{.ArchivedAt}}</div>
				<div><a href="{{.Address}}">View Original</a></div>
			</div>
	  </header>
	  <article>
			<small>{{.Error}}</small>
	  </article>
	</body>
</html>
	`, ModernNormalizeStyles(), Styles())
}
