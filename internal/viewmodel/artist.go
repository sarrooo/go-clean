package viewmodel

// swagger:parameters getArtistController
type GetArtistRequest struct {
	// The artist id.
	// Required: true
	// in:path
	ID string `json:"id" uri:"id" binding:"required"`
}

// swagger:response getArtistController
type GetArtistResponse struct {
	// in:body
	Body struct {
		// The artist id.
		// Required: true
		ID string `json:"id"`

		// The artist name.
		// Required: true
		Name string `json:"name"`
	} `json:"body"`
}
