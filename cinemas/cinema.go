/**
 * Author: Gaston Siffert
 * Created Date: 2017-11-08 20:11:02
 * Last Modified: 2017-11-11 11:11:24
 * Modified By: Gaston Siffert
 */

package cinemas

import (
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Cinema struct {
	Name       string
	Address    string
	DetailLink string
	// Location Location
}

// type Location struct {
// 	Street string
// 	Zip    string
// 	City   string
// }

// // MUST be called with a selection targeting a cinema division,
// // from the summary page.
// // Page: http://www.allocine.fr/salle/{page_area}
// // Division: "div.theaterblock.j_entity_container"
func (c *Cinema) fromSummary(s *goquery.Selection) {
	a := s.Find("h2").Find("a")

	link, exist := a.Attr("href")
	if exist {
		link = path.Join(base, link)
	}
	name := a.Text()
	address := s.Find("p.lighten").Text()

	// Assign the value
	*c = Cinema{
		Name:       strings.TrimSpace(name),
		Address:    strings.TrimSpace(address),
		DetailLink: link,
	}
}
