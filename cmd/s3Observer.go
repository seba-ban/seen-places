/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jackc/pgx/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/seba-ban/seen-places/broker"
	commonflags "github.com/seba-ban/seen-places/commonFlags"
	"github.com/seba-ban/seen-places/formats"
	events "github.com/seba-ban/seen-places/protogo"
	"github.com/seba-ban/seen-places/queries"
	"github.com/seba-ban/seen-places/utils"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

func getFormatToQueueMapping() map[formats.KnownFormatType]string {
	formatQueueMap := make(map[formats.KnownFormatType]string)
	formatQueueMap[formats.GarminFitFormatType] = commonflags.ChannelNames.GarminPointsExtract
	formatQueueMap[formats.GoProVideoFormatType] = commonflags.ChannelNames.GoproPointsExtract
	return formatQueueMap
}

func downloadFileFromS3(
	sess *session.Session, objKey string,
) (string, error) {

	f, err := os.CreateTemp(commonflags.LocalPaths.TmpDir, "*."+objKey)
	if err != nil {
		return "", err
	}
	defer f.Close()

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(commonflags.S3ConnConfig.RawFilesBucket),
		Key:    aws.String(objKey),
	})

	return f.Name(), err
}

func listObjectsInS3Bucket(
	sess *session.Session, bucket string,
) ([]string, error) {

	svc := s3.New(sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}

	objKeys := make([]string, 0)
	for _, item := range resp.Contents {
		objKeys = append(objKeys, *item.Key)
	}

	return objKeys, nil
}

func deleteObjectFromS3(
	sess *session.Session, bucket string, objKey string,
) error {
	svc := s3.New(sess)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objKey),
	})
	return err
}

func copyFile(sourcePath string, targetPath string) error {
	src, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

type FileInfo struct {
	FilePath     string
	OriginalName string
	Size         int64
	Format       formats.KnownFormatType
}

func processObject(sess *session.Session, objKey string) (*FileInfo, error) {
	log.Printf("Processing object: %s", objKey)
	localPath, err := downloadFileFromS3(sess, objKey)
	if localPath != "" {
		defer os.Remove(localPath)
	}

	if err != nil {
		return nil, err
	}

	format := formats.CheckFormatType(localPath)
	if format == formats.UnknownFormatType {
		return nil, fmt.Errorf("unknown format for object: %s", objKey)
	}

	targetPath, err := utils.PrepareTargetFolder(localPath, commonflags.LocalPaths.LocalStorageDir)
	if err != nil {
		return nil, err
	}

	fullTargetPath := path.Join(commonflags.LocalPaths.LocalStorageDir, targetPath)
	_, err = os.Stat(fullTargetPath)

	if os.IsNotExist(err) {
		log.Printf("Copying file to: %s", fullTargetPath)
		err = copyFile(localPath, fullTargetPath)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		log.Printf("File already exists: %s", fullTargetPath)
	}

	info, err := os.Stat(fullTargetPath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		FilePath:     targetPath,
		OriginalName: path.Base(objKey),
		Size:         info.Size(),
		Format:       format,
	}, nil
}

func saveFileInfoToDb(ctx context.Context, info *FileInfo, q *queries.Queries) error {
	_, err := q.CreateDataSource(ctx, &queries.CreateDataSourceParams{
		Filepath:         info.FilePath,
		OriginalFilename: info.OriginalName,
		Size:             info.Size,
		Type:             string(info.Format),
	})
	return err
}

func sendFileInfoToQueue(info *FileInfo, channel *amqp.Channel, queueName string) error {
	req := &events.ProcessFileRequest{Filepath: info.FilePath}

	data, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return channel.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			Body: data,
		})
}

func processFiles() error {
	formatQueue := getFormatToQueueMapping()

	// Connect to broker
	conn, err := commonflags.BrokerConnConfig.OpenTransport()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Open channel
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// Declare queues
	err = broker.DeclareQueue(commonflags.ChannelNames.GarminPointsExtract, channel)
	if err != nil {
		return err
	}
	err = broker.DeclareQueue(commonflags.ChannelNames.GoproPointsExtract, channel)
	if err != nil {
		return err
	}

	// Connect to the db
	ctx, dbConn, err := commonflags.DbConnConfig.OpenDbConnection()
	if err != nil {
		return err
	}
	defer dbConn.Close(*ctx)

	// Prepare queries
	q := queries.New(dbConn)

	// Get S3 session
	sess, err := commonflags.S3ConnConfig.GetS3Session()
	if err != nil {
		return err
	}

	objects, err := listObjectsInS3Bucket(sess, commonflags.S3ConnConfig.RawFilesBucket)
	if err != nil {
		return err
	}

	for _, objKey := range objects {
		info, err := processObject(sess, objKey)
		if err != nil {
			log.Printf("Error processing object: %s, %s", objKey, err)
			continue
		}

		queueName, ok := formatQueue[info.Format]
		if !ok {
			log.Printf("No queue for format: %s", info.Format)
			continue
		}

		_, err = q.GetDataSourceByFilePath(*ctx, info.FilePath)

		if err == pgx.ErrNoRows {
			log.Printf("Saving file info to db: %s", info.FilePath)
			err = saveFileInfoToDb(*ctx, info, q)
			if err != nil {
				log.Printf("Error saving file info to db: %s", err)
				continue
			}
		} else if err != nil {
			log.Printf("Error getting data source by file path: %s", err)
			continue
		}

		count, err := q.GetDataSourcePointsCount(*ctx, info.FilePath)
		if err != nil {
			log.Printf("Error getting data source points count: %s", err)
			continue
		}

		if count == 0 {
			log.Printf("Sending message to queue: %s", queueName)
			err = sendFileInfoToQueue(info, channel, queueName)
			if err != nil {
				log.Printf("Error sending message to queue: %s", err)
				continue
			}
		}

		log.Printf("Deleting object: %s", objKey)
		err = deleteObjectFromS3(sess, commonflags.S3ConnConfig.RawFilesBucket, objKey)
		if err != nil {
			log.Printf("Error deleting object: %s", err)
			continue
		}
	}
	return nil
}

// s3ObserverCmd represents the s3Observer command
var s3ObserverCmd = &cobra.Command{
	Use:   "s3Observer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		for {
			processFiles()
			log.Printf("Loop finished, sleeping before the next one")
			time.Sleep(2 * time.Minute)
		}
	},
}

func init() {
	commonflags.AddBrokerConfigFlags(s3ObserverCmd)
	commonflags.AddDbConfigFlags(s3ObserverCmd)
	commonflags.AddChannelNamesFlags(s3ObserverCmd)
	commonflags.AddLocalPathsFlags(s3ObserverCmd)
	commonflags.AddS3ConnConfigFlags(s3ObserverCmd)

	rootCmd.AddCommand(s3ObserverCmd)
}
