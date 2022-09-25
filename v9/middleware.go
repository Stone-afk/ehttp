package v9

type Middleware func(next HandleFunc) HandleFunc
