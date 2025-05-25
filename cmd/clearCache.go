/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/ekant01/caching-proxy/internal/proxy"
	"github.com/spf13/cobra"
)

// clearCacheCmd represents the clearCache command
var clearCacheCmd = &cobra.Command{
	Use:   "clear-cache",
	Short: "Clear the caching proxy cache",

	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Clearing cache...")
		proxy.ClearAll()
		log.Println("Cache cleared successfully.")
	},
}

func init() {
	rootCmd.AddCommand(clearCacheCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearCacheCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearCacheCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
