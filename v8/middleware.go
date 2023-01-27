package v8

type Middleware func(next HandleFunc) HandleFunc
