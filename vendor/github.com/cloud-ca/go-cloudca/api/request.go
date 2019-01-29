package api

//HTTP Methods
const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

//A request object
type CcaRequest struct {
	Method   string
	Endpoint string
	Body     []byte
	Options  map[string]string
}
