package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var (
	listen    string
	processor string
)

// LineProcessor is the interface that plugins must implement
type LineProcessor interface {
	ProcessLine(line string) (bool, string)
	IsResponseComplete(accumulatedOutput string) bool
}

// FactorioLineProcessor implements LineProcessor for Factorio server
type FactorioLineProcessor struct{}

func (f *FactorioLineProcessor) ProcessLine(line string) (bool, string) {
	// Implement Factorio-specific line processing here
	// Return true if the line should be included in the output, false otherwise
	return true, line
}

func (f *FactorioLineProcessor) IsResponseComplete(accumulatedOutput string) bool {
	// Implement logic to determine if the Factorio server response is complete
	// This could be based on certain keywords or patterns in the output
	return true // Placeholder implementation
}

// DefaultLineProcessor is a fallback processor
type DefaultLineProcessor struct{}

func (d *DefaultLineProcessor) ProcessLine(line string) (bool, string) {
	return true, line
}

func (d *DefaultLineProcessor) IsResponseComplete(accumulatedOutput string) bool {
	return true
}

func getLineProcessor(mode string) LineProcessor {
	switch mode {
	case "factorio":
		return &FactorioLineProcessor{}
	default:
		return &DefaultLineProcessor{}
	}
}

func splitDashArguments(args []string, argsLenAtDash int) ([]string, []string) {
	if argsLenAtDash < 0 {
		return args, []string{}
	}

	return args[:argsLenAtDash], args[argsLenAtDash:]
}

func processOutput(r io.Reader, processor LineProcessor) chan string {
	responseChan := make(chan string)
	outputChan := make(chan string)

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line) // Print every line to console

			if include, processedLine := processor.ProcessLine(line); include {
				outputChan <- processedLine + "\n"
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading output: %v", err)
		}
		close(outputChan)
	}()

	go func() {
		for line := range outputChan {
			select {
			case responseChan <- line:
				// Line sent to API response
			default:
				// If no one is receiving, just continue
			}
		}
		close(responseChan)
	}()

	return responseChan
}

func main() {
	root := cobra.Command{
		Use:  "stdserve",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.ArgsLenAtDash() <= 0 {
				return errors.New("no executable provided, please provide at least one executable to run by using \"stdserve -- <executable>\"")
			}

			_, executableArgs := splitDashArguments(args, cmd.ArgsLenAtDash())
			if len(executableArgs) == 0 {
				return errors.New("no executable provided, please provide at least one executable to run by using \"stdserve -- <executable>\"")
			}

			fmt.Println("Identified executable:", executableArgs[0], " with arguments:", executableArgs[1:])

			execCmd := exec.Command(executableArgs[0], executableArgs[1:]...)

			stdinPipe, err := execCmd.StdinPipe()
			if err != nil {
				return err
			}

			stdoutPipe, err := execCmd.StdoutPipe()
			if err != nil {
				return err
			}

			execCmd.Stderr = os.Stderr

			lineProcessor := getLineProcessor(processor)
			responseChan := processOutput(stdoutPipe, lineProcessor)

			e := echo.New()
			e.POST("/api/v1/stdin/execute", func(c echo.Context) error {
				type executeBody struct {
					Input string `json:"input"`
				}

				var body executeBody
				err := c.Bind(&body)
				if err != nil {
					return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
				}

				_, err = stdinPipe.Write([]byte(body.Input + "\n"))
				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
				}

				var accumulatedOutput strings.Builder
				timer := time.NewTimer(10 * time.Second)
				defer timer.Stop()

				for {
					select {
					case output, ok := <-responseChan:
						if !ok {
							// Channel closed, no more output
							return c.JSON(http.StatusOK, map[string]string{"output": accumulatedOutput.String()})
						}
						accumulatedOutput.WriteString(output)
						if lineProcessor.IsResponseComplete(accumulatedOutput.String()) {
							return c.JSON(http.StatusOK, map[string]string{"output": accumulatedOutput.String()})
						}
					case <-timer.C:
						return c.JSON(http.StatusRequestTimeout, map[string]string{"error": "timeout waiting for response"})
					}
				}
			})

			httpServer := &http.Server{
				Addr:    listen,
				Handler: e,
			}

			go func() {
				ln, err := net.Listen("tcp", listen)
				if err != nil {
					log.Fatalf("Error listening on %s: %v", listen, err)
				}

				fmt.Println("stdserve API Listening on", listen)

				err = httpServer.Serve(ln)
				if err != nil && err != http.ErrServerClosed {
					log.Printf("HTTP server error: %v", err)
				}
			}()

			err = execCmd.Start()
			if err != nil {
				return err
			}

			return execCmd.Wait()
		},
	}

	root.Flags().StringVarP(&listen, "listen", "l", "0.0.0.0:10080", "listen to interface & port")
	root.Flags().StringVarP(&processor, "processor", "p", "default", "processing mode (e.g., 'factorio', 'default')")

	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
