package routes

import (
	"log"
	"net/http"

	"github.com/fabianofski/equaly-backend/api/bucket"
	"github.com/fabianofski/equaly-backend/api/db"
	"github.com/labstack/echo/v4"
)

// HandlerGetProfilePicture godoc
//
//	@Summary		Get Profile Picture presigned Url
//	@Description	Retrieves presigned url of profile picture
//	@Tags			Expenses
//	@Param			expenseListId	path		string	true	"Expense List Id"
//	@Param			participantId	path		string	true	"Participant Id"
//	@Param			Authorization	header		string	true	"Bearer Token"
//	@Success		200				"Success"	string	"Profile Picture presigned URL"
//	@Failure		400				"Bad Request"
//	@Failure		500				"Internal Server Error"
//	@Router			/static/profile/{expenseListId}/{participantId} [get]
func HandlerGetProfilePicture(c echo.Context) error {
	userId := c.Get("userId").(string)
	log.Println("GET Profile Picture")

	expenseListId := c.Param("expenseListId")
	participantId := c.Param("participantId")

	if expenseListId == "" || participantId == "" {
		log.Println("400 Bad Request")
		return c.String(http.StatusBadRequest, "400 Bad Request")
	}

	authorized, err := db.IsUserAuthorized(expenseListId, userId)
	if err != nil {
		log.Println("500 Internal Server Error")
		log.Println(err)
		return c.String(http.StatusInternalServerError, "500 Internal Server Error")
	}

	if !authorized {
		log.Println("403 Forbidden")
		return c.String(http.StatusForbidden, "403 Forbidden")
	}

	url, err := bucket.CreatePresignedUrl("equaly", expenseListId+"/profile-"+participantId)
	if err != nil {
		log.Println("500 Internal Server Error")
		log.Println(err)
		return c.String(http.StatusInternalServerError, "500 Internal Server Error")
	}

	log.Println("200 Success")
	return c.String(http.StatusOK, url)
}

// HandlerPostProfilePicture godoc
//
//	@Summary		Upload Profile Picture of Participant
//	@Description	Upload Profile Picture of Participant
//	@Tags			Expenses
//	@Param			expenseListId	path		string	true	"Expense List Id"
//	@Param			participantId	path		string	true	"Participant Id"
//	@Param			image			formData	file	true	"Profile	Picture"
//	@Param			Authorization	header		string	true	"Bearer Token"
//	@Success		200				"Success"
//	@Failure		400				"Bad Request"
//	@Failure		500				"Internal Server Error"
//	@Router			/static/profile/{expenseListId}/{participantId} [Post]
func HandlerPostProfilePicture(c echo.Context) error {
	userId := c.Get("userId").(string)
	log.Println("POST Profile Picture")

	expenseListId := c.Param("expenseListId")
	participantId := c.Param("participantId")

	header, err := c.FormFile("image")
	if err != nil {
		log.Println("400 Bad Request")
		log.Println(err)
		return c.String(http.StatusBadRequest, "400 Bad Request")
	}

	src, err := header.Open()
	defer src.Close()
	if err != nil {
		log.Println("400 Bad Request")
		log.Println("Error opening file", err)
		return c.String(http.StatusBadRequest, "400 Bad Request")
	}

	if expenseListId == "" || participantId == "" {
		log.Println("400 Bad Request")
		return c.String(http.StatusBadRequest, "400 Bad Request")
	}

	isOwner, err := db.IsUserOwner(expenseListId, userId)
	if err != nil {
		log.Println("500 Internal Server Error")
		log.Println(err)
		return c.String(http.StatusInternalServerError, "500 Internal Server Error")
	}

	if !isOwner {
		log.Println("403 Forbidden")
		return c.String(http.StatusForbidden, "403 Forbidden")
	}

	err = bucket.UploadFile("equaly", expenseListId+"/profile-"+participantId, src, header.Size, header.Header.Get("Content-Type"))
	if err != nil {
		log.Println("500 Internal Server Error")
		log.Println(err)
		return c.String(http.StatusInternalServerError, "500 Internal Server Error")
	}

	log.Println("200 Success")
	return c.String(http.StatusOK, "200 Success")
}
