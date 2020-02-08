package template

func Styles() string {
	return `
		body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue", sans-serif;
			font-size: 1.3em;
		}
		
		img {
			max-width: 100%;
		}

		body > header {
			width: 80%;
			margin: 0 auto;
			margin-bottom: 10em;
			text-align: center;
		}

		body > header small, body > header .archived-at {
			color: gray;
		}

		body > header .archived-at {
			font-size: 50%;
		}

		body > header img {
			margin-bottom: 1em;
		}

		body > header figcaption {
			font-size: 80%;
		}

		body > article {
			width: 45%;
			margin: 0 auto;
		}

		a {
			text-decoration: none;
		}

		a:hover {
			text-decoration: underline;
		}
`
}

