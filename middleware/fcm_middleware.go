package middleware

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var MessagingClient *messaging.Client

func InitFirebase(c *gin.Context) {
	ctx := context.Background()

	opt := option.WithAuthCredentialsFile(option.AuthorizedUser, "firebase-service-account.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Gagal inisialisasi",
			},
		)

		return
	}

	MessagingClient, err = app.Messaging(ctx)
	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "Gagal inisialisasi",
			},
		)

		return
	}

	c.Set("fcm", MessagingClient)
	c.Next()
}
