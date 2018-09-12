package slack

import (
	api "github.com/nlopes/slack"
	"log"
	"strings"
)

// A Handler responds to an Message
type Handler interface {
	ServeMessage(*api.MessageEvent, *api.Client)
}

type ServeMux struct {
	m map[string]muxEntry
}

type muxEntry struct {
	h       Handler
	pattern string
}

func NewServeMux() *ServeMux { return new(ServeMux) }

var DefaultServeMux = &defaultServeMux

var defaultServeMux ServeMux

func (mux *ServeMux) match(text string) (h Handler, pattern string) {
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

func (mux *ServeMux) Handler(ev *api.MessageEvent, client *api.Client) (h Handler, pattern string) {
	log.Printf("[INFO] message is %s", ev.Msg.Text)
	h, pattern = mux.match(ev.Msg.Text)

	if h == nil {
		log.Printf("[INFO] not found pattern for %s", ev.Msg.Text)
		h, pattern = NothingHandler(), ""
	}

	return
}

func (mux *ServeMux) Handle(pattern string, handler Handler) {
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
		mux.m = make(map[string]muxEntry)
	}
	mux.m[pattern] = muxEntry{h: handler, pattern: pattern}
}

func (mux *ServeMux) ServeMessage(ev *api.MessageEvent, client *api.Client) {
	h, _ := mux.Handler(ev, client)
	h.ServeMessage(ev, client)
}

type HandlerFunc func(ev *api.MessageEvent, client *api.Client)

// ServeMessage calls f(w, r).
func (f HandlerFunc) ServeMessage(ev *api.MessageEvent, client *api.Client) {
	f(ev, client)
}

func Nothing(ev *api.MessageEvent, client *api.Client) {
	params := api.NewPostMessageParameters()
	if _, _, err := client.PostMessage(ev.Channel, "Command not found :innocent:", params); err != nil {
		log.Printf("[ERROR] failed to post message: %s", err)
	}
}

func NothingHandler() Handler {
	return HandlerFunc(Nothing)
}
