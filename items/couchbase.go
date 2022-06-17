package main

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

type ItemCouchbase struct {
	bucket  *gocb.Bucket
	cluster *gocb.Cluster
}

func NewCouchbase() ItemCouchbase {

	// cluster details
	endpoint := "127.0.0.1"
	bucketName := "items"
	username := "admin"
	password := "abc123"

	// Initialize the Connection
	cluster, err := gocb.Connect("couchbases://"+endpoint, gocb.ClusterOptions{

		Username: username,
		Password: password,
		SecurityConfig: gocb.SecurityConfig{
			TLSSkipVerify: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	bucket := cluster.Bucket(bucketName)

	err = bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		log.Fatal(err)
	}

	cluster.Query("CREATE PRIMARY INDEX ON `items`.`_default`.`_default`;", &gocb.QueryOptions{})

	return ItemCouchbase{bucket, cluster}
}

func (ic ItemCouchbase) Add(item Item) {

	// Get a reference to the default collection, required for older Couchbase server versions
	// col := bucket.DefaultCollection()

	col := ic.bucket.Scope("_default").Collection("_default")

	// Create and store a Document
	_, err := col.Upsert(string(item.Id),
		item, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (ic ItemCouchbase) Get(id int, user string) Item {
	results, err := ic.cluster.Query("SELECT items.* FROM items WHERE items.`username`=$1 AND `id` = $2;",
		&gocb.QueryOptions{PositionalParameters: []interface{}{user, id}})

	var h Item
	results.Next()
	err = results.Row(&h)
	fmt.Println(h.Id)
	if err != nil {
		panic(err)
	}

	return h
}

func (ic ItemCouchbase) GetAll(user string) []Item {
	//scope := ic.bucket.Scope("_default")
	results, err := ic.cluster.Query("SELECT items.* FROM items WHERE items.`username`=$1;",
		&gocb.QueryOptions{PositionalParameters: []interface{}{user}})
	// check query was successful
	if err != nil {
		panic(err)
	}

	var items []Item
	for results.Next() {
		var h Item
		err := results.Row(&h)
		fmt.Println(h.Id)
		if err != nil {
			panic(err)
		}
		items = append(items, h)
	}

	return items
}
