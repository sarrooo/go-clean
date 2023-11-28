package viewmodel

// swagger:parameters createArtistController
type CreateArtistRequest struct {
	// The artist name.
	// Required: true
	// in:body
	Name string `json:"name" binding:"required"`
}

// swagger:response createArtistController
type CreateArtistResponse struct {
	// in:body
	Body struct {
		// The artist id.
		// Required: true
		ID uint `json:"id"`
	} `json:"body"`
}

// swagger:parameters getArtistController
type GetArtistRequest struct {
	// The artist id.
	// Required: true
	// in:path
	ID uint `json:"id" uri:"id" binding:"required"`
}

// swagger:response getArtistController
type GetArtistResponse struct {
	// in:body
	Body struct {
		// The artist id.
		// Required: true
		ID uint `json:"id"`

		// The artist name.
		// Required: true
		Name string `json:"name"`
	} `json:"body"`
}

// swagger:parameters deleteArtistController
type DeleteArtistRequest struct {
	// The artist id.
	// Required: true
	// in:path
	ID uint `json:"id" uri:"id" binding:"required"`
}

// swagger:response deleteArtistController
type DeleteArtistResponse struct {
	// in:body
	Body struct {
		// The artist id.
		// Required: true
		ID string `json:"id"`
	} `json:"body"`
}
