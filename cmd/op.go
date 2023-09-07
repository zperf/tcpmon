package cmd

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var opCmd = &cobra.Command{
	Use:   "op",
	Short: "Operation commands",
}

var opBackupAllCmd = &cobra.Command{
	Use:     "backup [ADDR]",
	Short:   "Backup cluster",
	Example: " backup http://192.168.228.2:6789",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		baseUrl := args[0]
		timeout := viper.GetDuration("timeout")

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		backupUrl, err := url.JoinPath(baseUrl, "members")
		if err != nil {
			log.Fatal().Err(err).Msg("join url failed")
		}

		rsp, err := FetchJSON(ctx, backupUrl)
		if err != nil {
			log.Fatal().Err(err).Msg("Fetch API failed")
		}

		date := time.Now().Format(time.DateOnly)
		dir := "tcpmon-dump-" + date

		fmt.Println("#!/usr/bin/env bash")
		fmt.Println("set -x")
		fmt.Println("mkdir -p " + dir)
		fmt.Println("pushd " + dir + " || exit")
		for _, memberInfo := range rsp["members"].(map[string]any) {
			m := memberInfo.(map[string]any)
			fmt.Printf("curl -JfSsLO %s/backup\n", m["httpListen"])
		}
		fmt.Println("popd || exit")
		fmt.Printf("tar -czvf %s.tar.gz %s\n", dir, dir)
		fmt.Println("rm -rf " + dir)
	},
}

func init() {
	opBackupAllCmd.Flags().DurationP("timeout", "t", 4*time.Second, "HTTP timeout")
	rootCmd.AddCommand(opCmd)
	opCmd.AddCommand(opBackupAllCmd)
}
