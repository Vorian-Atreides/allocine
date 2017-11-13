/**
 * Author: Gaston Siffert
 * Created Date: 2017-11-09 18:11:50
 * Last Modified: 2017-11-09 18:11:10
 * Modified By: Gaston Siffert
 */

package cinemas

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

const (
	areaURL = "http://www.allocine.fr/salle/"
)

var (
	// we define Paris as a constant, because every links in the
	// Paris region and in Paris give a similar list of cinemas
	paris = Area{
		Link: "http://www.allocine.fr/salle/cinemas-pres-de-115755/",
		Name: "Paris",
	}
)

type Area struct {
	Name string
	Link string
}

func (c CinemaScraper) Areas() ([]Area, error) {
	doc, err := goquery.NewDocument(areaURL)
	if err != nil {
		return nil, err
	}

	areas := []Area{paris}
	doc.Find("div#region_120003").Find("a.underline").
		Each(func(i int, s *goquery.Selection) {
			link, exist := s.Attr("href")
			if !exist {
				return
			}
			link = fmt.Sprintf("%s%s", base, link)
			name := s.Find("span").Text()

			area := Area{Link: link, Name: name}
			areas = append(areas, area)
		})
	return areas, nil
}
