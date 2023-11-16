/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/seba-ban/seen-places/broker"
	commonflags "github.com/seba-ban/seen-places/commonFlags"
	events "github.com/seba-ban/seen-places/protogo"
	"github.com/seba-ban/seen-places/queries"
	"github.com/seba-ban/seen-places/utils"
	"github.com/spf13/cobra"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
	"github.com/twpayne/go-geom/encoding/wkbcommon"
	"google.golang.org/protobuf/proto"
)

// dbSaverCmd represents the dbSaver command
var dbSaverCmd = &cobra.Command{
	Use:   "dbSaver",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		conn, err := commonflags.BrokerConnConfig.OpenTransport()
		utils.PanicOnError(err)
		defer conn.Close()

		channel, err := conn.Channel()
		utils.PanicOnError(err)
		defer channel.Close()

		err = broker.DeclareQueue(commonflags.ChannelNames.PointsSave, channel)
		utils.PanicOnError(err)

		ctx, dbConn, err := commonflags.DbConnConfig.OpenDbConnection()
		utils.PanicOnError(err)
		defer dbConn.Close(*ctx)

		q := queries.New(dbConn)

		err = serveQueue(channel, q)
		utils.PanicOnError(err)

		fmt.Println("serving queue")

		var forever chan struct{}
		<-forever
	},
}

func init() {
	commonflags.AddBrokerConfigFlags(dbSaverCmd)
	commonflags.AddDbConfigFlags(dbSaverCmd)
	commonflags.AddChannelNamesFlags(dbSaverCmd)

	rootCmd.AddCommand(dbSaverCmd)
}

func processDelivery(d *amqp.Delivery, q *queries.Queries) {
	points := &events.ExtractedFilePoints{}
	err := proto.Unmarshal(d.Body, points)

	if err != nil {
		log.Printf("Error unmarshalling: %s", err.Error())
		d.Ack(false)
		return
	}

	log.Printf("received %d points", len(points.GetPoints()))

	if len(points.GetPoints()) == 0 {
		log.Println("no points received")
		d.Ack(false)
		return
	}

	// TODO: would be easier and more efficient to use pgx copy from directly...
	params := make([]*queries.CreatePointsParams, len(points.GetPoints()))

	for i, p := range points.GetPoints() {
		point := geom.NewPoint(geom.XY)
		point.SetSRID(4326)
		point.SetCoords([]float64{float64(p.Longitude), float64(p.Latitude)})
		raw, err := wkb.Marshal(point, wkbcommon.NDR)
		if err != nil {
			log.Printf("Error marshalling: %s", err.Error())
			d.Ack(false)
			return
		}

		// eh
		t, err := time.Parse(time.RFC3339, p.GetTimestamp())
		if err != nil {
			log.Printf("Error parsing timestamp: %s", err.Error())
			d.Ack(false)
			return
		}

		params[i] = &queries.CreatePointsParams{
			Geom:               raw,
			Visited:            pgtype.Timestamptz{Time: t, Valid: true},
			DataSourceFilepath: p.GetFilepath(),
		}
	}

	ctx := context.Background()
	saved, err := q.CreatePoints(ctx, params)

	if err != nil {
		log.Printf("Error saving: %s", err.Error())
		// TODO: check if it's a transient error or something like fk contraint violation
		d.Ack(false)
		return
	}

	log.Printf("Saved %d points", saved)
	d.Ack(false)
}

func serveQueue(ch *amqp.Channel, q *queries.Queries) error {
	msgs, err := ch.Consume(
		commonflags.ChannelNames.PointsSave, // queue
		"",                                  // consumer
		false,                               // auto-ack
		false,                               // exclusive
		false,                               // no-local
		false,                               // no-wait
		nil,                                 // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			processDelivery(&d, q)
		}
	}()

	return nil

}
