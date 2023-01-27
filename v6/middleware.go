package v6

type Middleware func(next HandleFunc) HandleFunc
