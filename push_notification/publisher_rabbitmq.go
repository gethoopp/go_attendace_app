package push_notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func Publisher_mssg(c *gin.Context, msg string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  http.StatusBadGateway,
			"message": "cannot connect to rabbitmq",
		})
		return
	}

	defer conn.Close()

	channel, err := conn.Channel()

	if err != nil {
		panic(err)
	}

	queue, err := channel.QueueDeclare(
		"Halo user",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	err = channel.Publish(
		"",
		"Halo user",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Succesfully publis message to" + queue.Name,
	})

}
