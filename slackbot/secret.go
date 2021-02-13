package slackbot

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
)

var (
	secretClient *secretmanager.Client
	projectNum   string
)

func getSecret(ctx context.Context, name string) (string, error) {
	resp, err := secretClient.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectNum, name),
	})
	if err != nil {
		return "", fmt.Errorf("access latest secret version of %s: %w", name, err)
	}
	return string(resp.Payload.GetData()), nil
}
