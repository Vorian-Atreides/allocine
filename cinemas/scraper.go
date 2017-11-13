/**
 * Author: Gaston Siffert
 * Created Date: 2017-11-08 20:11:40
 * Last Modified: 2017-11-11 11:11:12
 * Modified By: Gaston Siffert
 */

package cinemas

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/Vorian-Atreides/allocine/utils"
)

const (
	base = "http://www.allocine.fr"
)

type CinemaScraper struct {
	NumberOfWorkers uint
}

func (c CinemaScraper) PerAreaAndPage(a Area, page int) ([]Cinema, error) {
	url := fmt.Sprintf("%s?page=%d", a.Link, page)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	cinemas := []Cinema{}
	doc.Find("div.theaterblock.j_entity_container").
		Each(func(i int, s *goquery.Selection) {
			// Instantiate and fill the structure
			cinema := Cinema{}
			cinema.fromSummary(s)
			// Save the newly parsed cinema
			cinemas = append(cinemas, cinema)
		})
	return cinemas, nil
}

func (c CinemaScraper) PerArea(a Area) ([]Cinema, error) {
	nbPages, err := c.perAreaPageCount(a)
	if err != nil {
		return nil, err
	}

	toWorkers := make(chan cinemaAsyncRequest, c.NumberOfWorkers)
	fromWorkers := make(chan cinemaAsyncResponse, nbPages)
	defer func() {
		close(toWorkers)
		close(fromWorkers)
	}()

	// Run the workers
	for i := uint(0); i < c.NumberOfWorkers; i++ {
		go c.cinemaWorker(toWorkers, fromWorkers)
	}

	// Feed the workers
	for i := 0; i < nbPages; i++ {
		toWorkers <- cinemaAsyncRequest{area: a, page: i + 1}
	}

	// Aggregate the results
	cinemas := []Cinema{}
	for i := 0; i < nbPages; i++ {
		response := <-fromWorkers
		err = utils.ErrorConcat(err, response.err)
		if response.err == nil {
			cinemas = append(cinemas, response.cinemas...)
		}
	}
	return cinemas, err
}

// We could also make a loop of request until we receive a 302,
// but I don't consider it as efficient
func (c CinemaScraper) perAreaPageCount(a Area) (int, error) {
	doc, err := goquery.NewDocument(a.Link)
	if err != nil {
		return 0, err
	}

	// Get the last button in the "pages menu"
	number := doc.Find("table.centeringtable").Find("li.navcenterdata").
		Children().Last().Text()
	return strconv.Atoi(number)
}
