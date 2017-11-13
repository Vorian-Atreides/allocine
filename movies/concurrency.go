/**
 * Author: Gaston Siffert
 * Created Date: 2017-11-05 14:11:14
 * Last Modified: 2017-11-10 22:11:48
 * Modified By: Gaston Siffert
 */

package movies

//
// Summarize
//

type movieAsyncRequest struct {
	page int
}

type movieAsyncResponse struct {
	movies []Movie
	err    error
}

// const (
// 	defau
// )

func (m MovieScraper) summarizeWorker(input <-chan movieAsyncRequest,
	output chan<- movieAsyncResponse) {

	for {
		// Receive a task
		request := <-input
		// Handle the close(input)
		if request == (movieAsyncRequest{}) {
			return
		}

		movies, err := m.Summarize(request.page)
		output <- movieAsyncResponse{movies: movies, err: err}
	}
}
