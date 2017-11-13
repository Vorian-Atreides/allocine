/**
 * Author: Gaston Siffert
 * Created Date: 2017-11-09 09:11:30
 * Last Modified: 2017-11-09 09:11:15
 * Modified By: Gaston Siffert
 */

package utils

import (
	"fmt"
)

func ErrorConcat(dest error, src error) error {
	if src == nil {
		return dest
	}
	if dest == nil {
		return src
	}
	return fmt.Errorf("%s\n%s", dest.Error(), src.Error())
}
