/**
 * Author: Gaston Siffert
 * Created Date: 2017-11-09 18:11:50
 * Last Modified: 2017-11-09 18:11:13
 * Modified By: Gaston Siffert
 */

package utils

func UintMin(a uint, b uint) uint {
	if a < b {
		return a
	}
	return b
}
