package ES

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"main/pkg/util"
)

func InitCilent(host string, port int) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("%s:%d", host, port),
		},
		// Username: "elastic",
		// Password: "",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		util.Log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		util.Log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		util.Log.Fatalf("Error: %s", res.String())
	}

	util.Log.Println("成功连接到 Elasticsearch ")
	return es
}
