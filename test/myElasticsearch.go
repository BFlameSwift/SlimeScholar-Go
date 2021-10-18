package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/olivere/elastic"
)

var client *elastic.Client

var host = "http://82.156.217.192:9200"

