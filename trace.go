package ngtelgcp

import (
	"context"
	"fmt"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcpdetector "go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// GetTracerOpts gives us options used for the tracer and resource for tracing with Google Cloud.
func GetTracerOpts() ([]sdktrace.TracerProviderOption, []resource.Option, error) {
	exporter, err := texporter.New(
		texporter.WithTraceClientOptions([]option.ClientOption{option.WithTelemetryDisabled()}),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	return []sdktrace.TracerProviderOption{sdktrace.WithBatcher(exporter)},
		[]resource.Option{resource.WithDetectors(gcpdetector.NewDetector())},
		nil
}

// GetCloudTracePath gives us the path identifier of our current trace, enabling us to connect it in logs for example.
// It returns an empty string if the trace ID isn't valid or Google Cloud project ID could not be found.
func GetTracePath(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if projectID == "" {
		creds, _ := google.FindDefaultCredentials(ctx)
		projectID = creds.ProjectID
	}

	if projectID == "" || !sc.TraceID().IsValid() {
		return ""
	}

	tracePath := fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID())

	return tracePath
}
