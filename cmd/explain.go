/*
Copyright © 2025 Roman Acheampong romanacheampong4002@gmail.com
*/
package cmd

import (
	"github.com/r0mankeys/palace/tui"
	"github.com/spf13/cobra"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "Learn about the Method of Loci and how Palace works",
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.RunExplainTUI()
	},
}

func init() {
	rootCmd.AddCommand(explainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
