package slack

import (
	api "github.com/nlopes/slack"
	"log"
	"net/http"
)

// A InteractiveActionHandler responds to an Message
type InteractiveActionHandler interface {
	ServeInteractiveAction(*api.AttachmentActionCallback, http.ResponseWriter)
}

type ServeInteractiveActionMux struct {
	m          map[string]interactiveActionMuxEntry
	callbackID string
}

type interactiveActionMuxEntry struct {
	h       InteractiveActionHandler
	pattern string
}

func NewServeInteractiveActionMux(callbackID string) *ServeInteractiveActionMux {
	mux := new(ServeInteractiveActionMux)
	mux.callbackID = callbackID
	return mux
}

var DefaultServeInteractiveActionMux = &defaultServeInteractiveActionMux

var defaultServeInteractiveActionMux ServeInteractiveActionMux

func (mux *ServeInteractiveActionMux) match(name string) (h InteractiveActionHandler, pattern string) {
	v, ok := mux.m[name]
	if ok {
		return v.h, v.pattern
	}

	return
}

func (mux *ServeInteractiveActionMux) InteractiveActionHandler(callback *api.AttachmentActionCallback, w http.ResponseWriter) (h InteractiveActionHandler, pattern string) {
	action := callback.Actions[0]
	log.Printf("[INFO] callback is %s %d", action.Name, len(mux.m))
	h, pattern = mux.match(action.Name)

	if h == nil {
		log.Printf("[INFO] not found pattern for %s", action.Name)
		h, pattern = InteractiveActionNotFoundHandler(), ""
	}

	return
}

func (mux *ServeInteractiveActionMux) Handle(pattern string, handler InteractiveActionHandler) {
	if pattern == "" {
		panic("slack: invalid pattern")
	}
	if handler == nil {
		panic("slack: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("slack: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]interactiveActionMuxEntry)
	}
	mux.m[pattern] = interactiveActionMuxEntry{h: handler, pattern: pattern}
	log.Printf("[INFO] add pattern %s (%d)", pattern, len(mux.m))
}

func (mux *ServeInteractiveActionMux) ServeInteractiveAction(callback *api.AttachmentActionCallback, w http.ResponseWriter) {
	h, _ := mux.InteractiveActionHandler(callback, w)
	h.ServeInteractiveAction(callback, w)
}

type InteractiveActionHandlerFunc func(callback *api.AttachmentActionCallback, w http.ResponseWriter)

// ServeInteractiveAction calls f(w, r).
func (f InteractiveActionHandlerFunc) ServeInteractiveAction(callback *api.AttachmentActionCallback, w http.ResponseWriter) {
	f(callback, w)
}

func InteractiveActionNotFound(callback *api.AttachmentActionCallback, w http.ResponseWriter) {
	log.Printf("[ERROR] ]Invalid callback was submitteds")
	w.WriteHeader(http.StatusInternalServerError)
}

func InteractiveActionNotFoundHandler() InteractiveActionHandler {
	return InteractiveActionHandlerFunc(InteractiveActionNotFound)
}
