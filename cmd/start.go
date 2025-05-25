/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ekant01/caching-proxy/internal/proxy"
	"github.com/spf13/cobra"
)

var (
	port   int
	origin string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the caching proxy server",
	Long: `Start the caching proxy server that listens on a specified port and forwards requests to an origin server.
This command initializes the server and begins processing incoming requests, caching responses as needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
		err := proxy.StartServer(port, origin)
		if err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		} else {
			fmt.Printf("Caching proxy server started on port %d, forwarding to origin %s\n", port, origin)
		}

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the caching proxy server on")
	startCmd.Flags().StringVarP(&origin, "origin", "o", "http://localhost:8000", "Origin server URL to forward requests to")
	startCmd.MarkFlagRequired("origin") // Ensure origin is provided
}
