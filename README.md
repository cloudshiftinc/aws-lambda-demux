[![tests][1]][2]
[![Go Reference][3]][4]
[![GoCard][5]][6]
[![codecov][7]][8]
[![Apache V2 License](https://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/cloudshiftinc/aws-lambda-demux/blob/main/LICENSE)

[1]: https://github.com/cloudshiftinc/aws-lambda-demux/workflows/tests/badge.svg
[2]: https://github.com/cloudshiftinc/aws-lambda-demux/actions?query=workflow%3Atests
[3]: https://pkg.go.dev/badge/github.com/cloudshiftinc/aws-lambda-demux.svg
[4]: https://pkg.go.dev/github.com/cloudshiftinc/aws-lambda-demux
[5]: https://goreportcard.com/badge/github.com/cloudshiftinc/aws-lambda-demux
[6]: https://goreportcard.com/report/github.com/cloudshiftinc/aws-lambda-demux
[7]: https://codecov.io/gh/cloudshiftinc/aws-lambda-demux/branch/main/graph/badge.svg
[8]: https://codecov.io/gh/cloudshiftinc/aws-lambda-demux

Library to help Go developers handle multiple types of events (de-multiplexing) in AWS Lambda functions.

# Getting Started

The primary function of this library is to create events of a specific type and dispatch those to appropriate handlers.

To do so the demuxer is configured with `Factory` and `Handler` instances.

Factories are responsible for determining the type of the event (based off the incoming JSON) and creating an instance of that event.

Handlers are responsible for, well, handling that event.  Handlers are as used in [aws-lambda-go ](https://github.com/aws/aws-lambda-go), with the restriction
of having a signature of `func(context.Context, *eventType) (*responseType, error)`.  `eventType` and `responseType` can be any struct with the appropriate json tags to map from the event JSON.

A minimal usage showing a lambda that handles REST API request and Websocket lifecycle events:

```go
// main.go
package main

import (
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/cloudshiftinc/aws-lambda-demux/demux"
)

func main() {

  cfg := &demux.Cfg{
    Factories: []demux.Factory{
      func(ctx *demux.EventContext) any {
        if demux.HasAttribute(ctx.Event, "connectionId") {
          return &events.APIGatewayWebsocketProxyRequest{}
        }
        return &events.APIGatewayProxyRequest{}
      },
    },
    Handlers: []any{
      func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (
        *events.APIGatewayProxyResponse, error) {
        // TODO - your code here to handle websocket event
        return &events.APIGatewayProxyResponse{}, nil
      },
      func(ctx context.Context, event *events.APIGatewayProxyRequest) (
        *events.APIGatewayProxyResponse,
        error) {
        // TODO - your code here to handle HTTP/REST event
        return &events.APIGatewayProxyResponse{}, nil
      },
    },
  }

  lambda.Start(demux.NewHandler(cfg))
}

```

This library is not limited to event types in aws-lambda-go; any event type (including your own custom ones) that as appropriate JSON mappings can be used.

