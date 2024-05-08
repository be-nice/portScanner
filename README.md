# Port Scanner

Port Scanner is a simple command-line tool written in Go for scanning open ports on a specified IP address within a given range of port numbers.

## Overview

The Port Scanner utilizes Go's `net` package to establish TCP connections to each port within the specified range on the target IP address. It then reports back which ports are open and, if known, identifies the service commonly associated with each open port.

## Usage

To use the Port Scanner, follow these steps:

1. Clone the repository to your local machine:

2. Navigate to the project directory:


3. Build the project:

4. Run the scanner with the desired IP address and port range. If no port range is provided, it defaults to scanning 1000 most common open ports according to nmap.  
If only one port is provided, then it scans that one port only:
```bash
./portscan <IP_address> [<start_port> [<end_port>]]
```

Example:
```bash
./portscan 192.168.1.1 1 1000
```

This command will scan ports 1 through 1000 on the IP address `192.168.1.1`.

## Port configuration testing

Port scanner supports sqlite database to store expected port configuration for an ip.  
Schema is. ip | port | status  
where both ip and port are desired strings, and status is either "open" or "closed", depending on if you expect the port to be open or closed.  
Running the tool with ip and -t flag tests only the ports defined in the database.

```bash
./portscan 127.0.0.1 -t
```

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
