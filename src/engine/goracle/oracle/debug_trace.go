// +build trace

package oracle

import _ "log"

const CTrace = false

// prints with log.Printf the C-call trace
func ctrace(name string, args ...interface{}) {
	_log.Printf("CTRACE %s(%v)", name, args)
}
