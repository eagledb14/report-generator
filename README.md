# Form-Generator

Form Generator for creating reports quickly for Cyber Security Alerting purposes.

## Installation

https://github.com/eagledb14/shodan-form.git

Make sure download golang at least 1.23.2 or later

Make is not required, though it has a makefile included to help with building and distribution

## Usage

In the form-generator file 
```make```
runs the program and should open the program in a new browser window.

If the browser window does not open, the terminal shows the port on localhost that it is running on. It will run the first available port on startup.


```make build``` creates a zip file that has all the required files to run on windows and linux.

Windows Defender does not like the .exe so an exception needs to be made for it to run.
