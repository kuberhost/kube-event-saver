package pkg

import (
	"k8s.io/api/core/v1"
)

type IncomingEvent struct {
	Action string
	Event  *v1.Event
}
