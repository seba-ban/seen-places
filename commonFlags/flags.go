package commonflags

import (
	"os"

	"github.com/seba-ban/seen-places/broker"
	"github.com/seba-ban/seen-places/storage"
	"github.com/spf13/cobra"
)

type ChannelNamesStruct struct {
	PointsSave          string
	GoproPointsExtract  string
	GarminPointsExtract string
}

type LocalPathsStruct struct {
	LocalStorageDir string
	TmpDir          string
}

var DbConnConfig = &storage.DbConnectionConfig{}
var BrokerConnConfig = &broker.BrokerConnectionConfig{}
var ChannelNames = &ChannelNamesStruct{}
var LocalPaths = &LocalPathsStruct{}
var S3ConnConfig = &storage.S3ConnConfigStruct{}

func AddDbConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&DbConnConfig.Host, "db-host", "", "localhost", "Database host",
	)
	cmd.PersistentFlags().IntVarP(
		&DbConnConfig.Port, "db-port", "", 5432, "Database port",
	)
	cmd.PersistentFlags().StringVarP(
		&DbConnConfig.DbName, "db-name", "", "", "Database name",
	)
	cmd.PersistentFlags().StringVarP(
		&DbConnConfig.User, "db-user", "", "", "Database user",
	)
	cmd.PersistentFlags().StringVarP(
		&DbConnConfig.Password, "db-password", "", "", "Database password",
	)
	cmd.PersistentFlags().StringVarP(
		&DbConnConfig.Sslmode, "db-sslmode", "", "require", "Database sslmode, cf. https://www.postgresql.org/docs/current/libpq-ssl.html#LIBPQ-SSL-PROTECTION",
	)

	cmd.MarkPersistentFlagRequired("db-host")
	cmd.MarkPersistentFlagRequired("db-name")
	cmd.MarkPersistentFlagRequired("db-user")
	cmd.MarkPersistentFlagRequired("db-password")
}

func AddBrokerConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&BrokerConnConfig.Protocol, "broker-protocol", "", broker.AMQP_PROTOCOL, "BROKER protocol",
	)
	cmd.PersistentFlags().StringVarP(
		&BrokerConnConfig.Host, "broker-host", "", "", "BROKER host",
	)
	cmd.PersistentFlags().IntVarP(
		&BrokerConnConfig.Port, "broker-port", "", 5672, "AMQP port",
	)
	cmd.PersistentFlags().StringVarP(
		&BrokerConnConfig.User, "broker-user", "", "", "AMQP user",
	)
	cmd.PersistentFlags().StringVarP(
		&BrokerConnConfig.Password, "broker-password", "", "", "AMQP password",
	)
	cmd.PersistentFlags().StringVarP(
		&BrokerConnConfig.Vhost, "broker-vhost", "", "/", "AMQP vhost",
	)

	cmd.MarkPersistentFlagRequired("broker-host")
	cmd.MarkPersistentFlagRequired("broker-user")
	cmd.MarkPersistentFlagRequired("broker-password")
}

func AddChannelNamesFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&ChannelNames.PointsSave, "points-save-channel", "", "points", "Channel name for saving points",
	)
	cmd.PersistentFlags().StringVarP(
		&ChannelNames.GoproPointsExtract, "gopro-points-extract-channel", "", "gopro", "Channel name for extracting points from GoPro videos",
	)
	cmd.PersistentFlags().StringVarP(
		&ChannelNames.GarminPointsExtract, "garmin-points-extract-channel", "", "garmin", "Channel name for extracting points from Garmin FIT files",
	)
}

func AddLocalPathsFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&LocalPaths.LocalStorageDir, "local-storage-dir", "", "", "Local storage directory",
	)
	// TODO: check if dir exists
	cmd.PersistentFlags().StringVarP(
		&LocalPaths.TmpDir, "tmp-dir", "", os.TempDir(), "Temporary directory",
	)

	cmd.MarkPersistentFlagRequired("local-storage-dir")
}

func AddS3ConnConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(
		&S3ConnConfig.Region, "s3-region", "", "default", "S3 region",
	)
	cmd.PersistentFlags().StringVarP(
		&S3ConnConfig.RawFilesBucket, "s3-raw-files-bucket", "", "", "S3 raw files bucket",
	)
	cmd.PersistentFlags().StringVarP(
		&S3ConnConfig.AccessKey, "s3-access-key", "", "", "S3 access key",
	)
	cmd.PersistentFlags().StringVarP(
		&S3ConnConfig.SecretKey, "s3-secret-key", "", "", "S3 secret key",
	)
	cmd.PersistentFlags().StringVarP(
		&S3ConnConfig.EndpointUrl, "s3-endpoint-url", "", "", "S3 endpoint URL",
	)

	cmd.MarkPersistentFlagRequired("s3-raw-files-bucket")
	cmd.MarkPersistentFlagRequired("s3-access-key")
	cmd.MarkPersistentFlagRequired("s3-secret-key")
	cmd.MarkPersistentFlagRequired("s3-endpoint-url")
}
