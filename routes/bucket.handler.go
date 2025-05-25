package routes

import (
	"net/http"

	"github.com/fabianofski/equaly-backend/api/bucket"
	"github.com/labstack/echo/v4"
)

// HandlerGetProfilePicture godoc
//
//	@Summary		Get Profile Picture presigned Url
//	@Description	Retrieves presigned url of profile picture
//	@Tags			Expenses
//	@Param			Authorization	header		string	true	"Bearer Token"
//	@Success		200				"Success"	string	"Profile Picture presigned URL"
//	@Failure		400				"Bad Request"
//	@Failure		500				"Internal Server Error"
//	@Router			/static/profile [get]
func HandlerGetProfilePicture(c echo.Context) error {

	url, err := bucket.CreatePresignedUrl("equaly", "giphy.gif")
	if err != nil {
		return c.String(http.StatusInternalServerError, "500 Internal Server Error")
	}

	return c.String(http.StatusOK, url)
}
