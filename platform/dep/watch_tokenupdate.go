package dep

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micromdm/micromdm/platform/config"
	"github.com/micromdm/micromdm/platform/pubsub"
)

func (svc *DEPService) watchTokenUpdates(pubsub pubsub.Subscriber) error {
	tokenAdded, err := pubsub.Subscribe(context.TODO(), "list-token-events", config.DEPTokenTopic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-tokenAdded:
				var token config.DEPToken
				if err := json.Unmarshal(event.Message, &token); err != nil {
					fmt.Printf("unmarshalling tokenAdded to token: %v\n", err)
					continue
				}

				client, err := token.Client()
				if err != nil {
					fmt.Printf("creating new DEP client: %v\n", err)
					continue
				}

				svc.mtx.Lock()
				svc.client = client
				svc.mtx.Unlock()
			}
		}
	}()

	return nil
}
