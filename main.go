/*
* Serial Data Logger - Open Source Contribution by Shuhain
*
* This program reads data from a serial USB port and logs it to a CSV file with timestamps
* that include milliseconds. Users can save or discard the data using keyboard shortcuts.
* This project follows high coding standards and is structured in a way that could
* comply with ISO 9001 principles for quality and continuous improvement.
*
* Keyboard Shortcuts:
*  - Alt+C: Save the CSV file and exit.
*  - Ctrl+X: Exit without saving (the CSV file will be discarded).
 */

package main

import (
	"encoding/csv" // For writing CSV files
	"fmt"          // For formatted I/O
	"io/ioutil"    // For reading and writing files
	"log"          // For logging errors
	"os"           // For interacting with the operating system (file and signal management)
	"strings"      // For string manipulations
	"time"         // For timestamps with millisecond precision

	"github.com/tarm/serial" // Package for handling serial port communication
	"golang.org/x/term"      // For reading raw terminal input (used to detect key combinations)
)

// Configuration file to store the last used COM port
const portConfigFile = "portconfig.txt"

// Global variables for managing the CSV file and data writer
var writer *csv.Writer
var file *os.File
var exitWithoutSaving = false // Flag to indicate if the program should exit without saving

/*
* saveComPort saves the provided COM port to a file for future use.
* This allows the user to avoid re-entering the COM port on subsequent runs.
 */
func saveComPort(comPort string) error {
	err := ioutil.WriteFile(portConfigFile, []byte(comPort), 0644)
	if err != nil {
		return err
	}
	return nil
}

/*
* loadComPort loads the previously saved COM port from the configuration file.
* Returns the COM port string if found, or an error if the file doesn't exist.
 */
func loadComPort() (string, error) {
	data, err := ioutil.ReadFile(portConfigFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

/*
* confirmComPort checks if a COM port was saved previously and asks the user
* if they want to use the same COM port. If the user declines, it prompts for a new port.
 */
func confirmComPort() string {
	comPort, err := loadComPort()
	if err == nil && comPort != "" {
		// Ask the user if they want to use the previous COM port
		fmt.Printf("Previously saved COM port: %s\n", comPort)
		fmt.Print("Do you want to use this COM port? (Y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" {
			return comPort
		}
	}

	// Prompt for a new COM port if no port is saved or the user declines
	fmt.Print("Enter the COM port (e.g., COM3): ")
	fmt.Scanln(&comPort)
	saveComPort(comPort) // Save the new COM port for future use
	return comPort
}

/*
* createCSVFile generates a new CSV file with the current date and time in its filename.
* This function returns the file and a CSV writer object to write data into the file.
 */
func createCSVFile() (*os.File, *csv.Writer, error) {
	// Generate a timestamp for the file name with the format YYYYMMDD_HHMMSS
	currentTime := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("data_%s.csv", currentTime) // Example: data_20240908_123456.csv

	// Create the CSV file for logging serial data
	file, err := os.Create(fileName)
	if err != nil {
		return nil, nil, err
	}

	// Initialize the CSV writer to write rows of data
	writer := csv.NewWriter(file)
	return file, writer, nil
}

/*
* checkKeys listens for keypresses in raw mode to detect Alt+C and Ctrl+X.
* - Alt+C: Save data and exit.
* - Ctrl+X: Exit without saving the CSV file.
 */
func checkKeys() {
	for {
		key := readChar() // Read raw terminal input

		if key == "\x18" { // ASCII value of Ctrl+X is 24
			// Exit without saving
			fmt.Println("\nCtrl+X pressed. Exiting without saving.")
			exitWithoutSaving = true
			if file != nil {
				file.Close()           // Close the file
				os.Remove(file.Name()) // Delete the unsaved CSV file
			}
			os.Exit(0) // Terminate the program
		} else if key == "\x1b" { // Escape character (Alt key detection)
			// Detect if Alt+C is pressed (Alt sends ESC followed by the character)
			nextKey := readChar()
			if nextKey == "c" || nextKey == "C" {
				// Save data and exit
				fmt.Println("\nAlt+C pressed. Saving data and exiting.")
				if writer != nil && file != nil {
					writer.Flush() // Ensure all buffered data is written to the file
					file.Close()   // Safely close the CSV file
				}
				os.Exit(0) // Terminate the program
			}
		}
	}
}

/*
* readChar reads a single character from the terminal in raw mode.
* This function is used to detect key combinations like Ctrl+X and Alt+C.
 */
func readChar() string {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(fd, oldState) // Ensure terminal is restored to normal mode after input

	var buf [1]byte
	os.Stdin.Read(buf[:])
	return string(buf[:])
}

/*
* main is the entry point of the program. It initializes the serial communication,
* prompts the user for the COM port (or confirms a saved port), starts logging serial
* data to a CSV file, and monitors for keypresses to save or discard the data.
 */
func main() {
	// Confirm and get the COM port from the user
	comPort := confirmComPort()

	// Configure the serial port with the default baud rate of 115200
	c := &serial.Config{Name: comPort, Baud: 115200} // Adjust to your device's settings
	s, err := serial.OpenPort(c)                     // Open the serial port
	if err != nil {
		log.Fatalf("Failed to open COM port: %v", err)
	}
	defer s.Close() // Ensure the port is closed when the program exits

	// Create a CSV file for logging the serial data
	file, writer, err = createCSVFile()
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close() // Ensure the file is closed when the program exits

	// Buffer to hold the serial data read from the device
	buf := make([]byte, 128)

	// Start a goroutine to listen for Alt+C and Ctrl+X keypresses
	go checkKeys()

	// Main loop: Continuously read data from the serial port and log it to the CSV file
	fmt.Println("Reading serial data. Press Alt+C to save and exit, or Ctrl+X to exit without saving.")
	for {
		// Read the data from the serial port into the buffer
		n, err := s.Read(buf)
		if err != nil {
			log.Fatalf("Failed to read from serial port: %v", err)
		}

		// Convert the received bytes into a string and display it in the terminal
		data := string(buf[:n])
		fmt.Printf("Data received: %s\n", data)

		// Write the timestamp (with millisecond precision) and the data to the CSV file
		currentTime := time.Now().Format("2006-01-02 15:04:05.000")
		err = writer.Write([]string{currentTime, data})
		if err != nil {
			log.Fatalf("Failed to write to CSV: %v", err)
		}
		writer.Flush() // Ensure the data is written to the CSV file
	}
}
