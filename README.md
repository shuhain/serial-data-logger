
# Serial Data Logger

This project is an open-source serial data logger that reads data from a serial USB port and logs it to a CSV file with timestamps (including milliseconds). The user can save or discard the data using keyboard shortcuts. It is designed to be used for open-source contributions.

## Features
- Logs serial data in real-time.
- Saves the data to a CSV file with a timestamp (including milliseconds).
- Automatically detects previously used COM port and prompts for confirmation.
- Uses keyboard shortcuts for control:
  - `Alt+C` to save the data and exit.
  - `Ctrl+X` to discard the data and exit without saving.

## Installation and Setup

### Prerequisites
- [Go (Golang)](https://golang.org/doc/install) 1.16 or higher
- Git installed on your machine ([Download Git](https://git-scm.com/))
- Serial USB device or any microcontroller with serial communication

### Installing the Project

1. **Clone the Repository**:
   Open your terminal and clone the repository using the following command:
   
   ```bash
   git clone https://github.com/yourusername/yourrepositoryname.git
   cd yourrepositoryname
   ```

2. **Install the Required Go Dependencies**:
   You need to install the `tarm/serial` package for serial communication and the `golang.org/x/term` package for raw terminal input. Run the following commands:
   
   ```bash
   go get github.com/tarm/serial
   go get golang.org/x/term
   ```

3. **Run the Application**:
   Now, run the Go program:
   
   ```bash
   go run main.go
   ```

   If this is the first time running, it will ask you for the COM port (e.g., `COM3` or `/dev/ttyUSB0`). The port will be saved for future use and prompted for confirmation on subsequent runs.

### Usage

1. **Reading Serial Data**:
   - Once the program starts, it will begin reading data from the serial port and logging it to the terminal.
   - It will also log the data to a CSV file with a timestamp (in the format `YYYY-MM-DD HH:MM:SS.mmm` where `mmm` is milliseconds).
   
2. **Keyboard Shortcuts**:
   - Press **Alt+C** to save the CSV file and exit the program.
   - Press **Ctrl+X** to exit the program without saving the CSV file (the file will be deleted).

### Example CSV Output
Below is an example of what the output in the CSV file looks like:

```csv
Timestamp,Data
2024-09-08 13:45:23.123,Data received from the serial device
2024-09-08 13:45:23.234,Another data entry
```

### Software Requirements
- **Go (Golang)** 1.16 or higher
- **Serial Device** (e.g., Arduino, Raspberry Pi with UART)
- **Operating System**: Windows, macOS, or Linux

### License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contribution
Contributions are welcome! If you find any issues or want to add new features, feel free to fork the repository and submit a pull request.

## Author
**Shuhain** - [GitHub](https://github.com/yourusername)
