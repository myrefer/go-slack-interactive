package slack

import (
	api "github.com/nlopes/slack"
	"log"
	"strings"
)

// A InteractiveActionHandler responds to an Message
type InteractiveActionHandler interface {
	ServeInteractiveAction(*api.MessageEvent, *api.Client)
}

type ServeInteractiveActionMux struct {
	m map[string]interactiveActionMuxEntry
}

type interactiveActionMuxEntry struct {
	h       InteractiveActionHandler
	pattern string
}

func NewServeInteractiveActionMux() *ServeInteractiveActionMux { return new(ServeInteractiveActionMux) }

var DefaultServeInteractiveActionMux = &defaultServeInteractiveActionMux

var defaultServeInteractiveActionMux ServeInteractiveActionMux

func (mux *ServeInteractiveActionMux) match(text string) (h InteractiveActionHandler, pattern string) {
	// Parse message
	cmd := strings.Split(strings.TrimSpace(text), " ")[1:]
	if len(cmd) == 0 {
		return
	}

	v, ok := mux.m[cmd[0]]
	if ok {
		return v.h, v.pattern
	}

	return
}

func (mux *ServeInteractiveActionMux) InteractiveActionHandler(ev *api.MessageEvent, client *api.Client) (h InteractiveActionHandler, pattern string) {
	log.Printf("[INFO] message is %s", ev.Msg.Text)
	h, pattern = mux.match(ev.Msg.Text)

	if h == nil {
		log.Printf("[INFO] not found pattern for %s", ev.Msg.Text)
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
}

func (mux *ServeInteractiveActionMux) ServeInteractiveAction(ev *api.MessageEvent, client *api.Client) {
	h, _ := mux.InteractiveActionHandler(ev, client)
	h.ServeInteractiveAction(ev, client)
}

type InteractiveActionHandlerFunc func(ev *api.MessageEvent, client *api.Client)

// ServeInteractiveAction calls f(w, r).
func (f InteractiveActionHandlerFunc) ServeInteractiveAction(ev *api.MessageEvent, client *api.Client) {
	f(ev, client)
}

func InteractiveActionNotFound(ev *api.MessageEvent, client *api.Client) {
	params := api.NewPostMessageParameters()
	if _, _, err := client.PostMessage(ev.Channel, "Command not found :innocent:", params); err != nil {
		log.Printf("[ERROR] failed to post message: %s", err)
	}
}

func InteractiveActionNotFoundHandler() InteractiveActionHandler {
	return InteractiveActionHandlerFunc(InteractiveActionNotFound)
}
