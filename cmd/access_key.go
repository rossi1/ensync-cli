package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ensync-cli/internal/api"
	"github.com/ensync-cli/internal/domain"
	"github.com/spf13/cobra"
)

func newAccessKeyCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access-key",
		Short: "Manage access keys",
	}

	cmd.AddCommand(
		newAccessKeyListCmd(client),
		newAccessKeyCreateCmd(client),
		newAccessKeyVerifyCmd(client),
		newAccessKeyPermissionsCmd(client),
	)

	return cmd
}

func newAccessKeyListCmd(client *api.Client) *cobra.Command {
	var pageIndex int
	var limit int
	var order string
	var orderBy string
	var accessKey string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List access keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := &api.ListParams{
				PageIndex: pageIndex,
				Limit:     limit,
				Order:     order,
				OrderBy:   orderBy,
				Filter: map[string]string{
					"accessKey": accessKey,
				},
			}

			keys, err := client.ListAccessKeys(context.Background(), params)
			if err != nil {
				return fmt.Errorf("failed to list access keys: %w", err)
			}

			return printJSON(cmd.OutOrStdout(), keys)
		},
	}

	cmd.Flags().IntVar(&pageIndex, "page", 0, "Page index")
	cmd.Flags().IntVar(&limit, "limit", 10, "Number of items per page")
	cmd.Flags().StringVar(&order, "order", "DESC", "Sort order (ASC/DESC)")
	cmd.Flags().StringVar(&orderBy, "order-by", "createdAt", "Field to order by")
	cmd.Flags().StringVar(&accessKey, "key", "", "Filter by access key")

	return cmd
}

func newAccessKeyCreateCmd(client *api.Client) *cobra.Command {
	var accessFile string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new access key",
		RunE: func(cmd *cobra.Command, args []string) error {
			var accessKey domain.AccessKey
			if accessFile != "" {
				data, err := os.ReadFile(accessFile)
				if err != nil {
					return fmt.Errorf("failed to read access file: %w", err)
				}

				if err := json.Unmarshal(data, &accessKey); err != nil {
					return fmt.Errorf("failed to parse access JSON: %w", err)
				}
			}

			createdKey, err := client.CreateAccessKey(context.Background(), &accessKey)
			if err != nil {
				return fmt.Errorf("failed to create access key: %w", err)
			}

			return printJSON(cmd.OutOrStdout(), createdKey)
		},
	}

	cmd.Flags().StringVar(&accessFile, "file", "", "JSON file containing access key configuration")

	return cmd
}

func newAccessKeyVerifyCmd(client *api.Client) *cobra.Command {
	var accessKey string

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify an access key",
		RunE: func(cmd *cobra.Command, args []string) error {
			if accessKey == "" {
				return fmt.Errorf("access key is required")
			}

			valid, err := client.VerifyAccessKey(context.Background(), accessKey)
			if err != nil {
				return fmt.Errorf("failed to verify access key: %w", err)
			}

			result := struct {
				Valid bool   `json:"valid"`
				Key   string `json:"key"`
			}{
				Valid: valid,
				Key:   accessKey,
			}

			return printJSON(cmd.OutOrStdout(), result)
		},
	}

	cmd.Flags().StringVar(&accessKey, "key", "", "Access key to verify")
	cmd.MarkFlagRequired("key")

	return cmd
}

func newAccessKeyPermissionsCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions",
		Short: "Manage access key permissions",
	}

	cmd.AddCommand(
		newAccessKeyGetPermissionsCmd(client),
		newAccessKeySetPermissionsCmd(client),
	)

	return cmd
}

func newAccessKeyGetPermissionsCmd(client *api.Client) *cobra.Command {
	var accessKey string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get access key permissions",
		RunE: func(cmd *cobra.Command, args []string) error {
			if accessKey == "" {
				return fmt.Errorf("access key is required")
			}

			permissions, err := client.GetAccessKeyPermissions(context.Background(), accessKey)
			if err != nil {
				return fmt.Errorf("failed to get permissions: %w", err)
			}

			return printJSON(cmd.OutOrStdout(), permissions)
		},
	}

	cmd.Flags().StringVar(&accessKey, "key", "", "Access key")
	cmd.MarkFlagRequired("key")

	return cmd
}

func newAccessKeySetPermissionsCmd(client *api.Client) *cobra.Command {
	var accessKey string
	var permissionsFile string

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set access key permissions",
		RunE: func(cmd *cobra.Command, args []string) error {
			if accessKey == "" {
				return fmt.Errorf("access key is required")
			}

			if permissionsFile == "" {
				return fmt.Errorf("permissions file is required")
			}

			data, err := os.ReadFile(permissionsFile)
			if err != nil {
				return fmt.Errorf("failed to read permissions file: %w", err)
			}

			var permissions domain.AccessKeyPermissions
			if err := json.Unmarshal(data, &permissions); err != nil {
				return fmt.Errorf("failed to parse permissions JSON: %w", err)
			}

			err = client.SetAccessKeyPermissions(context.Background(), accessKey, &permissions)
			if err != nil {
				return fmt.Errorf("failed to set permissions: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Permissions updated successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&accessKey, "key", "", "Access key")
	cmd.Flags().StringVar(&permissionsFile, "file", "", "JSON file containing permissions")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("file")

	return cmd
}