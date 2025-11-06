/*
Copyright © 2025 Roman Acheampong romanacheampong4002@gmail.com 
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "palace",
	Short: "Palace is a personal memory training tool that helps users build and strengthen their memory through the method of loci",
	Long:`
	Welcome to Palace 🏰  
Your personal memory training companion.

Data Directory: ~/.palace  
Database: ~/.palace/palace.db  

Available commands:
  create   Create a new memory palace
  list     View your existing palaces
  train    Start a training session
  stats    View your memory stats
  export   Export your palace data
  help     Show command help`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.palace.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


