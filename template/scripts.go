package template

func Scripts() string {
	return `
document.querySelectorAll("img").forEach(img => {
  img.removeAttribute("srcset")
  img.removeAttribute("width")
  img.removeAttribute("height")
})
`
}
