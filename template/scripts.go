package template

func Scripts() string {
	return `
document.addEventListener("DOMContentLoaded", () => {
	document.querySelectorAll("img").forEach(img => {
		if (img.naturalWidth === 0) {
			img.remove()
		} else {
			img.removeAttribute("srcset")
			img.removeAttribute("width")
			img.removeAttribute("height")
		}
	})

	document.querySelectorAll("iframe").forEach(iframe => {
		const iframeSrc = iframe.getAttribute("src")
		if (iframeSrc.startsWith("//")) iframe.setAttribute("src", "https:" + iframeSrc)
	})
})
`
}
