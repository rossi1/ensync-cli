package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ensync-cli/internal/api"
	"github.com/ensync-cli/internal/domain"
)

func newEventCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event",
		Short: "Manage events",
	}

	cmd.AddCommand(
		newEventListCmd(client),
		newEventCreateCmd(client),
		newEventUpdateCmd(client),
	)

	return cmd
}

func newEventListCmd(client *api.Client) *cobra.Command {
	var pageIndex int
	var limit int
	var order string
	var orderBy string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List events",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := &api.ListParams{
				PageIndex: pageIndex,
				Limit:     limit,
				Order:     order,
				OrderBy:   orderBy,
			}

			events, err := client.ListEvents(context.Background(), params)
			if err != nil {
				return fmt.Errorf("failed to list events: %w", err)
			}

			return printJSON(cmd.OutOrStdout(), events)
		},
	}

	cmd.Flags().IntVar(&pageIndex, "page", 0, "Page index")
	cmd.Flags().IntVar(&limit, "limit", 10, "Number of items per page")
	cmd.Flags().StringVar(&order, "order", "DESC", "Sort order (ASC/DESC)")
	cmd.Flags().StringVar(&orderBy, "order-by", "createdAt", "Field to order by")

	return cmd
}

func newEventCreateCmd(client *api.Client) *cobra.Command {
	var name string
	var payloadFile string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new event definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("name is required")
			}

			var payload map[string]string
			if payloadFile != "" {
				// Read payload from file
				data, err := os.ReadFile(payloadFile)
				if err != nil {
					return fmt.Errorf("failed to read payload file: %w", err)
				}

				if err := json.Unmarshal(data, &payload); err != nil {
					return fmt.Errorf("failed to parse payload JSON: %w", err)
				}
			}

			event := &domain.Event{
				Name:    name,
				Payload: payload,
			}

			ctx := context.Background()
			err := client.CreateEvent(ctx, event)
			if err != nil {
				return fmt.Errorf("failed to create event: %w", err)
			}

			fmt.Printf("Event '%s' created successfully\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Event name")
	cmd.Flags().StringVar(&payloadFile, "payload-file", "", "Path to JSON file containing payload")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newEventUpdateCmd(client *api.Client) *cobra.Command {
	var id string
	var name string
	var payloadFile string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing event definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == "" {
				return fmt.Errorf("id is required")
			}

			var payload map[string]string
			if payloadFile != "" {
				// Read payload from file
				data, err := os.ReadFile(payloadFile)
				if err != nil {
					return fmt.Errorf("failed to read payload file: %w", err)
				}

				if err := json.Unmarshal(data, &payload); err != nil {
					return fmt.Errorf("failed to parse payload JSON: %w", err)
				}
			}

			event := &domain.Event{
				ID:      id,
				Name:    name,
				Payload: payload,
			}

			ctx := context.Background()
			err := client.UpdateEvent(ctx, event)
			if err != nil {
				return fmt.Errorf("failed to update event: %w", err)
			}

			fmt.Printf("Event '%s' updated successfully\n", id)
			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Event ID")
	cmd.Flags().StringVar(&name, "name", "", "New event name")
	cmd.Flags().StringVar(&payloadFile, "payload-file", "", "Path to JSON file containing new payload")
	cmd.MarkFlagRequired("id")

	return cmd
}
