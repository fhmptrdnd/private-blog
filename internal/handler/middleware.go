package handler

import (
	"fmt"
	"net/http"
)

// middleware, function yang nerima function dan return function
// contoh higher-order function
type Middleware func(http.HandlerFunc) http.HandlerFunc

// withlogging, tambahin logging ke handler
func WithLogging(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log request
		fmt.Printf("[LOG] %s %s\n", r.Method, r.URL.Path)
		handler(w, r)
	}
}

// withmethodcheck, cek method http yang diizinkan
// ini contoh closure: function capture variable allowedmethod
func WithMethodCheck(allowedMethod string) Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != allowedMethod {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			handler(w, r)
		}
	}
}

// withpanicrecovery, tangkep panic biar ga crash
func WithPanicRecovery(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				fmt.Printf("[PANIC] %v\n", err)
			}
		}()
		handler(w, r)
	}
}

// chain, gabungin beberapa middleware jadi satu
// contoh function composition
// implementasi rekursif: base case + recursive case
func Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	// base case: kalo ga ada middleware lagi, return handler
	if len(middlewares) == 0 {
		return handler
	}
	// recursive case: apply middleware terakhir, lalu rekursi ke sisanya
	// middlewares diterapkan dari kanan ke kiri (reverse order)
	last := middlewares[len(middlewares)-1]
	rest := middlewares[:len(middlewares)-1]
	return Chain(last(handler), rest...)
}

// withusercontext, bikin closure yang capture user id handling
func WithUserContext(handler func(w http.ResponseWriter, r *http.Request, userID string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getOrCreateUserID(w, r) // Captured in closure
		handler(w, r, userID)
	}
}

type HandlerFunc = http.HandlerFunc

// composehandlers, gabungin beberapa handler jadi satu pipeline
func ComposeHandlers(handlers ...HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, h := range handlers {
			h(w, r)
		}
	}
}
