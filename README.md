# wifipass-go🚀  
A simple command-line tool written in Golang to retrieve WiFi credentials and generate QR codes for easy sharing.  

#### Note: Tested on Linux ONLY.

## Features ✅  
- 📶 **Show current WiFi password**  
- 📜 **List previously connected networks**  
- 🔑 **Retrieve a specific connection's password**  
- 📷 **Generate and display a QR code for sharing WiFi**  
- 💾 **Save the QR code to a file**  
- 🩺 **Run a diagnostic check for necessary dependencies**  

#### Note: Make sure go1.23.0 or later version is already installed on you system.

## Installation ⚡  
```
git clone https://github.com/joybiswas007/wifipass-go.git
cd wifipass-go
make deps
make build
sudo mv wifipass /usr/local/bin/
wifipass --help
wifipass --doctor (Check if all required dependencies are installed)
```

## Usage
```
Usage:
  -connection string
        Show password for a specific WiFi connection (run as sudo)
  -doctor
        Run diagnostics to check WiFi settings, configurations, and required packages
  -list
        List previously connected WiFi networks (run as sudo)
  -qr
        Generate QR code for the current WiFi password and display it
  -qr-path string
        Specify path to save the QR code
  -save-qr
        Save the generated QR code to a file instead of displaying it
  -version
        Display version and exit
```
## Note: Some commands require sudo access.

## Contributing 🤝
Contributions are welcome! Feel free to open issues or submit pull requests.

## License 📜
This project is licensed under the MIT License.

## Similar Projects & Inspiration 🌟
This project was built independently but shares similar goals with other tools like:

[wifi-password by sdushantha](https://github.com/sdushantha/wifi-password)<br/>
[wifi-password by rauchg](https://github.com/rauchg/wifi-password)
