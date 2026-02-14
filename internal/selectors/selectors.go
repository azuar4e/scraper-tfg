package selectors

import (
	"strings"
)

func Selector(url string) (string, string) {
	// añadir las webs que se deseen y una para pruebas local
	if strings.Contains(url, "amazon") {
		return "#title", "span.a-price span.a-offscreen"
	} else if strings.Contains(url, "pccomponentes") {
		return "#pdp-title", "#pdp-price-current-integer"
	} else {
		return "", ""
	}
}
