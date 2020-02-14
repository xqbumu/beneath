package http

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/beneath-core/control/entity"
	"github.com/beneath-core/pkg/log"
	"github.com/beneath-core/pkg/secrettoken"
	"github.com/beneath-core/pkg/ws"
	"github.com/beneath-core/gateway"
	"github.com/beneath-core/gateway/subscriptions"
)

// wsServer implements ws.Server
type wsServer struct{}

// KeepAlive implements ws.Server
func (s wsServer) KeepAlive(numClients int, elapsed time.Duration) {
	// log state
	log.S.Infow(
		"ws keepalive",
		"clients", numClients,
		"elapsed", elapsed,
	)
}

// InitClient implements ws.Server
func (s wsServer) InitClient(client *ws.Client, payload map[string]interface{}) error {
	// get secret
	secret := entity.Secret(&entity.AnonymousSecret{})
	tokenStr, ok := payload["secret"].(string)
	if ok {
		// parse token
		token, err := secrettoken.FromString(tokenStr)
		if err != nil {
			return fmt.Errorf("malformed secret")
		}

		// authenticate
		secret = entity.AuthenticateWithToken(client.Context, token)
		if secret == nil {
			return fmt.Errorf("couldn't authenticate secret")
		}
	}

	// set secret as state
	client.SetState(secret)

	// log
	s.logWithSecret(secret, "ws init client",
		"ip", client.GetRemoteAddr(),
	)

	return nil
}

// CloseClient implements ws.Server
func (s wsServer) CloseClient(client *ws.Client) {
	// the client state is the secret
	secret := client.GetState().(entity.Secret)

	// log
	s.logWithSecret(secret, "ws close client",
		"ip", client.GetRemoteAddr(),
		"reads", client.MessagesRead,
		"bytes_read", client.BytesRead,
		"writes", client.MessagesWritten,
		"bytes_written", client.BytesWritten,
		"elapsed", time.Since(client.StartTime),
		"error", client.Err,
	)
}

// StartQuery implements ws.Server
func (s wsServer) StartQuery(client *ws.Client, id ws.QueryID, payload map[string]interface{}) error {
	// get instance id
	instanceIDStr, ok := payload["instance_id"].(string)
	if !ok {
		return fmt.Errorf("payload must contain key 'instance_id'")
	}
	instanceID := uuid.FromStringOrNil(instanceIDStr)
	if instanceID == uuid.Nil {
		return fmt.Errorf("query is not a valid instance ID")
	}

	// get cursor
	cursorStr, ok := payload["cursor"].(string)
	if !ok {
		return fmt.Errorf("payload must contain key 'cursor'")
	}
	cursor, err := base64.StdEncoding.DecodeString(cursorStr)
	if err != nil {
		return fmt.Errorf("payload key 'cursor' must contain a base64-encoded cursor")
	}

	// the client state is the secret
	secret := client.GetState().(entity.Secret)

	// get cached stream
	stream := entity.FindCachedStreamByCurrentInstanceID(client.Context, instanceID)
	if stream == nil {
		return fmt.Errorf("stream not found")
	}

	// if batch, check committed
	if stream.Batch && !stream.Committed {
		return fmt.Errorf("batch has not yet been committed, and so can't be read")
	}

	// check permissions
	perms := secret.StreamPermissions(client.Context, stream.StreamID, stream.ProjectID, stream.Public, stream.External)
	if !perms.Read {
		return fmt.Errorf("token doesn't grant right to read this stream")
	}

	// check usage
	if !secret.IsAnonymous() {
		usage := gateway.Metrics.GetCurrentUsage(client.Context, secret.GetOwnerID())
		if !secret.CheckReadQuota(usage) {
			return fmt.Errorf("You have exhausted your quota")
		}
	}

	// get subscription channel
	cancel, err := gateway.Subscriptions.Subscribe(instanceID, cursor, func(msg subscriptions.Message) {
		// TODO: For now, just sending empty messages
		client.SendData(id, "")
	})
	if err != nil {
		return err
	}

	// set cancel as query state
	client.SetQueryState(id, cancel)

	return nil
}

// StopQuery implements ws.Server
func (s wsServer) StopQuery(client *ws.Client, id ws.QueryID) error {
	// the query state is the cancel func
	cancel := client.GetQueryState(id).(context.CancelFunc)
	cancel()

	// the client state is the secret
	secret := client.GetState().(entity.Secret)

	// log
	s.logWithSecret(secret, "ws stop query",
		"ip", client.GetRemoteAddr(),
		"id", id,
	)

	return nil
}

func (s wsServer) logWithSecret(sec entity.Secret, msg string, keysAndValues ...interface{}) {
	l := log.S

	if sec.IsUser() {
		l = l.With(
			"secret", sec.GetSecretID().String(),
			"user", sec.GetOwnerID().String(),
		)
	} else if sec.IsService() {
		l = l.With(
			"secret", sec.GetSecretID().String(),
			"service", sec.GetOwnerID().String(),
		)
	}

	l.Infow(msg, keysAndValues...)
}