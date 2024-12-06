package radius

import "sync"

type ServeMux struct {
	mu sync.RWMutex
	m  map[Code]Handler
}

func NewServeMux() *ServeMux {
	return new(ServeMux)
}

var DefaultServeMux = NewServeMux()

func (mux *ServeMux) Handle(code Code, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if handler == nil {
		panic("radius: nil handler")
	}
	if mux.m == nil {
		mux.m = make(map[Code]Handler)
	}
	if _, exist := mux.m[code]; exist {
		panic("radius: multiple registrations for " + code.String())
	}
	mux.m[code] = handler
}

func (mux *ServeMux) HandleFunc(code Code, handler func(ResponseWriter, *Request)) {
	mux.Handle(code, HandlerFunc(handler))
}

func (mux *ServeMux) match(code Code) Handler {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if h, ok := mux.m[code]; ok {
		return h
	}
	return nil
}

func (mux *ServeMux) ServeRADIUS(w ResponseWriter, r *Request) {
	h := mux.match(r.Code)
	h.ServeRADIUS(w, r)
}

func Handle(code Code, handler Handler) {
	DefaultServeMux.Handle(code, handler)
}

func HandleFunc(code Code, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(code, handler)
}
