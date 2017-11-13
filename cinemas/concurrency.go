/**
 * Author: Gaston Siffert
 * Created Date: 2017-11-09 18:11:35
 * Last Modified: 2017-11-10 22:11:33
 * Modified By: Gaston Siffert
 */

package cinemas

type cinemaAsyncRequest struct {
	area Area
	page int
}

type cinemaAsyncResponse struct {
	cinemas []Cinema
	err     error
}

func (c CinemaScraper) cinemaWorker(input <-chan cinemaAsyncRequest,
	output chan<- cinemaAsyncResponse) {

	for {
		// Receive a task
		request := <-input
		// Handle close(input)
		if request == (cinemaAsyncRequest{}) {
			return
		}

		cinemas, err := c.PerAreaAndPage(request.area, request.page)
		output <- cinemaAsyncResponse{cinemas: cinemas, err: err}
	}
}
