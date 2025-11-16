package middlewares

import (
	"fmt"
	"time"
	"net/http"
)

type ResponseTimer struct {
	http.ResponseWriter
	status int
}

func (rw *ResponseTimer) WriteHeader (code int){
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received Request in ResponseTime \n")
		start := time.Now()

		wrappedWriter := &ResponseTimer{
			ResponseWriter: w,
			status: http.StatusOK,
		}

		duration := time.Since(start)
		w.Header().Set("X-Response-Time", duration.String())
		next.ServeHTTP(wrappedWriter, r)

		duration = time.Since(start)
		fmt.Printf("[%d] %s %s - %v \n",
			wrappedWriter.status,
			r.Method,
			r.URL.Path,
			duration,
		)
		fmt.Println("Sent Request in ResponseTime")
	})
}
