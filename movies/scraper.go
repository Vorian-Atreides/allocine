/**
 * Author: Gaston Siffert
 * Created Date: 2017-11-04 21:11:01
 * Last Modified: 2017-11-11 00:11:08
 * Modified By: Gaston Siffert
 */

package movies

import (
	"fmt"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/Vorian-Atreides/allocine/utils"
)

const (
	summarizeURL = "http://www.allocine.fr/films/?page=%d"
)

// MovieScraper provide several methods to scrape the list of movies
// in allocine
type MovieScraper struct {
	NumberOfWorkers uint
}

func (m MovieScraper) Summarize(page int) ([]Movie, error) {
	url := fmt.Sprintf(summarizeURL, page)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	movies := []Movie{}
	doc.Find("div.card.card-entity.card-entity-list.cf").
		Each(func(i int, s *goquery.Selection) {
			// Instantiate and fill the movie structure
			movie := Movie{}
			movie.fromSummary(s)
			// Save the newly parsed movie
			movies = append(movies, movie)
		})
	return movies, nil
}

func (m MovieScraper) AllSummarize() ([]Movie, error) {
	nbPages, err := m.summarizePageCount()
	if err != nil {
		return nil, err
	}

	toWorkers := make(chan movieAsyncRequest, nbPages)
	fromWorkers := make(chan movieAsyncResponse, m.NumberOfWorkers)
	defer func() {
		close(toWorkers)
		close(fromWorkers)
	}()

	// Run the workers
	for i := uint(0); i < m.NumberOfWorkers; i++ {
		go m.summarizeWorker(toWorkers, fromWorkers)
	}

	// Feed the workers
	for i := 0; i < nbPages; i++ {
		toWorkers <- movieAsyncRequest{page: i + 1}
	}

	// Aggregate the results
	movies := []Movie{}
	for i := 0; i < nbPages; i++ {
		response := <-fromWorkers
		err = utils.ErrorConcat(err, response.err)
		if response.err == nil {
			movies = append(movies, response.movies...)
		}
	}
	return movies, err
}

// We could also make a loop of request until we receive a 302,
// but I don't consider it as efficient
func (m MovieScraper) summarizePageCount() (int, error) {
	url := fmt.Sprintf(summarizeURL, 1)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return 0, err
	}

	// Get the last button in the "pages menu"
	number := doc.Find("div.pagination-item-holder").
		Children().Last().Text()
	return strconv.Atoi(number)
}
