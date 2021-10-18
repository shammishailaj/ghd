package utils

import (
	"github.com/couchbase/gocb"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type CouchbaseDB struct {
	Conn    *gocb.Cluster
	ConnErr error
	Conf    *CouchbaseConfig
	Log     *log.Logger
}

func (c *CouchbaseDB) Connect() error {
	c.Conf = new(CouchbaseConfig)
	c.Conf.UserName = os.Getenv("WH_DB_USERNAME")
	c.Conf.Password = os.Getenv("WH_DB_PASSWORD")
	c.Conf.Port, _ = strconv.ParseInt(os.Getenv("WH_DB_PORT"), 10, 64)
	c.Conf.HostName = os.Getenv("WH_DB_HOSTNAME")
	c.Conf.Bucket = os.Getenv("WH_DB_BUCKET")
	c.Conf.BucketPass = os.Getenv("WH_DB_BUCKET_PASS")
	c.Conn, c.ConnErr = c.Conf.Connect()
	return c.ConnErr

}

func (c *CouchbaseDB) OpenBucket() (*gocb.Bucket, error) {
	// Open Bucket
	bucket, err := c.Conn.OpenBucket(c.Conf.Bucket, c.Conf.BucketPass)
	return bucket, err
}

func (c *CouchbaseDB) Upsert(key string, value ApiResponse) error {

	bucket, bucketErr := c.OpenBucket()
	if bucketErr != nil {
		c.Log.Printf("Error opening Bucket %s with password %s", c.Conf.Bucket, c.Conf.BucketPass)
		return bucketErr
	}

	_, casErr := bucket.Upsert(key, value, 0)
	return casErr
}
