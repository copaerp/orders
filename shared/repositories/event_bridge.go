package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	schedulersvc "github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/aws/smithy-go"
)

type EventBridgeClient struct {
	Scheduler *schedulersvc.Client
}

func NewEventBridgeClient() (*EventBridgeClient, error) {

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic("erro ao carregar config da AWS: " + err.Error())
	}
	scheduler := schedulersvc.NewFromConfig(cfg)

	return &EventBridgeClient{scheduler}, nil
}

func (c *EventBridgeClient) PutEvent(ctx context.Context, name string, timeout time.Duration, body map[string]any) error {
	scheduleTime := time.Now().Add(timeout).UTC().Format("2006-01-02T15:04:05")

	bodyStr, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body to JSON: %w", err)
	}
	strBody := string(bodyStr)

	input := &schedulersvc.CreateScheduleInput{
		Name:                       aws.String(name),
		ScheduleExpression:         aws.String(fmt.Sprintf("at(%s)", scheduleTime)),
		ScheduleExpressionTimezone: aws.String("UTC"),
		FlexibleTimeWindow: &types.FlexibleTimeWindow{
			Mode: types.FlexibleTimeWindowModeOff,
		},
		Target: &types.Target{
			Arn:     aws.String(os.Getenv("ot_arn")),
			RoleArn: aws.String(os.Getenv("role_arn")),
			Input:   aws.String(strBody),
		},
		GroupName: aws.String("order-lifecycle"),
	}

	_, err = c.Scheduler.CreateSchedule(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ConflictException" {
			updateInput := &schedulersvc.UpdateScheduleInput{
				Name:                       input.Name,
				ScheduleExpression:         input.ScheduleExpression,
				ScheduleExpressionTimezone: input.ScheduleExpressionTimezone,
				FlexibleTimeWindow:         input.FlexibleTimeWindow,
				Target:                     input.Target,
				GroupName:                  input.GroupName,
			}

			_, updateErr := c.Scheduler.UpdateSchedule(ctx, updateInput)
			if updateErr != nil {
				return fmt.Errorf("failed to update existing schedule: %w", updateErr)
			}
			return nil
		}
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	return nil
}
