package cmd

import (
	"context"
	"helmwatch/cmd/version"
	"helmwatch/internal/config"
	"helmwatch/internal/msg"
	"helmwatch/internal/tui"
	"helmwatch/internal/watcher"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sgtdi/fswatcher"
	"github.com/spf13/cobra"
)

const (
	templatesDir = "templates/"
)

func Execute() {
	rootCmd := new()

	rootCmd.AddCommand(version.New())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func new() *cobra.Command {
	var flags config.Flags

	cmd := &cobra.Command{
		Use:   "helmwatch",
		Short: "Interactive Helm diff watcher",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			configuration, err := config.Load(flags.ConfigPath)
			if err != nil {
				log.Fatal(err)
			}

			configuration.Flags = flags

			p := tea.NewProgram(tui.New(configuration), tea.WithAltScreen(), tea.WithContext(ctx))

			go func() {
				if err := watcher.Watch(ctx, func(_ fswatcher.WatchEvent) {
					p.Send(msg.FileChanged{})
				}, fswatcher.WithCooldown(500*time.Millisecond), fswatcher.WithIncRegex(flags.ValuesFile, templatesDir)); err != nil {
					log.Printf("failed to watch: %v", err)
				}
			}()

			if _, err := p.Run(); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&flags.Chart, "chart", ".", "Path to the chart or remote chart reference")
	cmd.Flags().StringVar(&flags.Version, "version", "", "Version of the chart")
	cmd.Flags().StringVar(&flags.ConfigPath, "config", "helmwatch.yaml", "Path to the helmwatch config")
	cmd.Flags().StringVar(&flags.ValuesFile, "values", "values.yaml", "Path to the values file")

	return cmd
}
