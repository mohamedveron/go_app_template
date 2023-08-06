package datastore

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Config struct {
	// URI: all configurations are expected to be part of the URI, e.g. read timeout
	URI string
	// ClientID can be the same as the application name, which is consuming this package
	ClientID    string
	dbname      string
	PingTimeout time.Duration

	ConnectTimeout         time.Duration
	HeartbeatInterval      time.Duration
	LocalThreshold         time.Duration
	MaxConnIdleTime        time.Duration
	MaxPoolSize            uint64
	MinPoolSize            uint64
	ServerSelectionTimeout time.Duration
}

// MongoDB - mongodb connection service. A single instance of this struct can be
// used to connect to only 1 database
type MongoDB struct {
	config *Config
	client *mongo.Client

	*mongo.Database
}

// Connect connects to the database, creating client
func (mdb *MongoDB) connect(
	ctx context.Context,
	opts *options.ClientOptions,
) error {
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to mongodb servers: %v", opts.Hosts)
	}

	mdb.client = mongoClient
	mdb.Database = mongoClient.Database(mdb.config.dbname)

	return nil
}

// Disconnect gracefully disconnects mongodb client
func (mdb *MongoDB) Disconnect(ctx context.Context) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, mdb.config.PingTimeout)
	defer cancel()
	if err := mdb.client.Disconnect(ctxWithTimeout); err != nil {
		return errors.Wrap(err, "failed to disconnect MongoDB")
	}

	return nil
}

// Collection returns mongodb collection by given name
func (mdb *MongoDB) Collection(name string) *mongo.Collection {
	return mdb.Database.Collection(name)
}

func (mdb *MongoDB) Client() *mongo.Client {
	return mdb.client
}

func ToObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func setDefaults(opts *options.ClientOptions) { //nolint
	if opts == nil {
		return
	}

	// SetConnectTimeout specifies a timeout that is used for creating connections to the server. If a custom Dialer is
	// specified through SetDialer, this option must not be used. This can be set through ApplyURI with the
	// "connectTimeoutMS" (e.g "connectTimeoutMS=30") option. If set to 0, no timeout will be used. The default is 30
	// seconds.
	if opts.Dialer == nil && (opts.ConnectTimeout == nil || *opts.ConnectTimeout <= 0) {
		opts.SetConnectTimeout(time.Second * 5) //nolint
	}

	// SetHeartbeatInterval specifies the amount of time to wait between periodic background server checks. This can also be
	// set through the "heartbeatIntervalMS" URI option (e.g. "heartbeatIntervalMS=10000"). The default is 10 seconds.
	if opts.HeartbeatInterval == nil || *opts.HeartbeatInterval < time.Second {
		opts.SetHeartbeatInterval(time.Second * 10) //nolint
	}

	// SetLocalThreshold specifies the width of the 'latency window': when choosing between multiple suitable servers for an
	// operation, this is the acceptable non-negative delta between shortest and longest average round-trip times. A server
	// within the latency window is selected randomly. This can also be set through the "localThresholdMS" URI option (e.g.
	// "localThresholdMS=15000"). The default is 15 milliseconds.
	if opts.LocalThreshold == nil || *opts.LocalThreshold < time.Millisecond {
		opts.SetLocalThreshold(time.Millisecond * 15) //nolint
	}

	// SetMaxConnIdleTime specifies the maximum amount of time that a connection will remain idle in a connection pool
	// before it is removed from the pool and closed. This can also be set through the "maxIdleTimeMS" URI option (e.g.
	// "maxIdleTimeMS=10000"). The default is 0, meaning a connection can remain unused indefinitely.
	if opts.MaxConnIdleTime == nil || *opts.MaxConnIdleTime <= 0 {
		opts.SetMaxConnIdleTime(time.Second * 60) //nolint
	}

	// SetMaxPoolSize specifies that maximum number of connections allowed in the driver's connection pool to each server.
	// Requests to a server will block if this maximum is reached. This can also be set through the "maxPoolSize" URI option
	// (e.g. "maxPoolSize=100"). If this is 0, maximum connection pool size is not limited. The default is 100.
	if opts.MaxPoolSize == nil || *opts.MaxPoolSize == 0 {
		opts.SetMaxPoolSize(100) //nolint
	}

	// SetMinPoolSize specifies the minimum number of connections allowed in the driver's connection pool to each server. If
	// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
	// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
	if opts.MinPoolSize != nil && *opts.MinPoolSize == 0 {
		opts.SetMinPoolSize(1)
	}

	// SetServerSelectionTimeout specifies how long the driver will wait to find an available, suitable server to execute an
	// operation. This can also be set through the "serverSelectionTimeoutMS" URI option (e.g.
	// "serverSelectionTimeoutMS=30000"). The default value is 30 seconds.
	if opts.ServerSelectionTimeout == nil || *opts.ServerSelectionTimeout < 0 {
		opts.SetServerSelectionTimeout(time.Second * 30) //nolint
	}

	// SetTimeout specifies the amount of time that a single operation run on this Client can execute before returning an error.
	// The deadline of any operation run through the Client will be honored above any Timeout set on the Client; Timeout will only
	// be honored if there is no deadline on the operation Context. Timeout can also be set through the "timeoutMS" URI option
	// (e.g. "timeoutMS=1000"). The default value is nil, meaning operations do not inherit a timeout from the Client.
	//
	// If any Timeout is set (even 0) on the Client, the values of MaxTime on operation options, TransactionOptions.MaxCommitTime and
	// SessionOptions.DefaultMaxCommitTime will be ignored. Setting Timeout and SocketTimeout or WriteConcern.wTimeout will result
	// in undefined behavior.
	//
	// NOTE(benjirewis): SetTimeout represents unstable, provisional API. The behavior of the driver when a Timeout is specified is
	// subject to change.
	if opts.Timeout == nil || *opts.Timeout < 0 {
		opts.SocketTimeout = nil
		opts.SetTimeout(time.Second * 15) //nolint
	}
}

func mongoOpts(cfg *Config) (*options.ClientOptions, string, error) {
	cs, err := connstring.ParseAndValidate(cfg.URI)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to parse connection URI")
	}

	opts := options.Client()
	opts.ApplyURI(cs.String())
	opts.SetAppName(cfg.ClientID)
	opts.SetHeartbeatInterval(cfg.HeartbeatInterval)
	opts.SetLocalThreshold(cfg.LocalThreshold)
	opts.SetMaxConnIdleTime(cfg.MaxConnIdleTime)
	opts.SetMaxPoolSize(cfg.MaxPoolSize)
	opts.SetMinPoolSize(cfg.MinPoolSize)
	opts.SetServerSelectionTimeout(cfg.ServerSelectionTimeout)
	opts.SetConnectTimeout(cfg.ConnectTimeout)

	setDefaults(opts)

	return opts, cs.Database, nil
}

func New(ctx context.Context, cfg *Config) (*MongoDB, error) {
	mdb := &MongoDB{
		config: cfg,
	}

	opts, dbname, err := mongoOpts(cfg)
	if err != nil {
		return nil, err
	}

	cfg.dbname = dbname

	if cfg.PingTimeout < time.Millisecond {
		cfg.PingTimeout = time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()

	err = mdb.connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = mdb.client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping MongoDB")
	}

	return mdb, nil
}
