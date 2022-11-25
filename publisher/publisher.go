package publisher

import (
	"apipublisher/transferobject"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

var publishers = make([]Publisher, 0)

func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}
	publisherKeys := strings.Split(os.Getenv("publishers"), ",")

	for _, value := range publisherKeys {
		if value == "AZURE" {
			publishers = append(publishers, &AzurePublisher{})
		}
		if value == "KONG" {
			publishers = append(publishers, &KongPublisher{})
		}
	}
}

type Publisher interface {
	PublishApi(api transferobject.Api) error

	PublishRoute(route transferobject.Route) error
}

type AzurePublisher struct{}

func (p *AzurePublisher) PublishApi(api transferobject.Api) error {
	fmt.Println("publishing api to azure")
	b, _ := json.Marshal(AzureApi{
		Name: api.Name,
		Type: api.Type,
	})
	fmt.Println(string(b))
	return nil
}

func (p *AzurePublisher) PublishRoute(route transferobject.Route) error {
	fmt.Println("publishing route to azure")
	b, _ := json.Marshal(AzureRoute{
		Name: route.Name,
		Path: route.Path,
	})
	fmt.Println(string(b))
	return nil
}

type KongPublisher struct{}

func (p *KongPublisher) PublishApi(api transferobject.Api) error {
	fmt.Println("publishing api to kong")
	b, _ := json.Marshal(KongApi{
		Name: api.Name,
		Type: api.Type,
	})
	fmt.Println(string(b))
	return nil
}

func (p *KongPublisher) PublishRoute(route transferobject.Route) error {
	fmt.Println("publishing route to kong")
	b, _ := json.Marshal(KongRoute{
		Name: route.Name,
		Path: route.Path,
	})
	fmt.Println(string(b))
	return nil
}

func Run(data string) error {
	var kafkaData transferobject.KafkaData
	err := json.Unmarshal([]byte(data), &kafkaData)

	if err == nil {
		fmt.Println(err)
	}

	for _, p := range publishers {
		if kafkaData.Event == "API" {
			var api transferobject.Api
			err := json.Unmarshal([]byte(kafkaData.Data), &api)

			if err != nil {
				fmt.Println(err)
			}
			err = p.PublishApi(api)
			if err != nil {
				return err
			}
		} else if kafkaData.Event == "ROUTE" {
			var route transferobject.Route
			err := json.Unmarshal([]byte(kafkaData.Data), &route)

			if err != nil {
				fmt.Println(err)
			}
			err = p.PublishRoute(route)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("not supported")
		}
	}
	return nil
}
