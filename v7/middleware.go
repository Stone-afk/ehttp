package v7

type Middleware func(next HandleFunc) HandleFunc
