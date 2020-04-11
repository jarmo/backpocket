package template

func Scripts() string {
	return `
document.querySelectorAll("img").forEach(img => {
  img.removeAttribute("srcset")
  img.removeAttribute("width")
  img.removeAttribute("height")
})

document.querySelectorAll("iframe").forEach(iframe => {
	const iframeSrc = iframe.getAttribute("src")
	if (iframeSrc.startsWith("//")) iframe.setAttribute("src", "https:" + iframeSrc)
})
`
}
