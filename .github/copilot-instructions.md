# Copilot Instructions for mister-mqtt

## Project Overview
This project provides an mqtt interface for the mister fpga platform. It reads the current state by reading the files /tmp/CORENAME, /tmp/ACTIVEGAME and /tmp/RBFNAME and publishes them to an mqtt broker. Is uses the homeassistant mqtt discovery format to make it easy to integrate with homeassistant. It uses fsnotify to watch for changes in the files and updates the mqtt broker when a change is detected.
This project is written in go. It uses /etc/hostname for getting the node_id for the homeassistant mqtt discovery.

## Configuration
The application is configured via flags. Only the mqtt broker address is required the default is localhost.

## Development Guidelines
- Use english for documentation, code, comments, user visible texts and logging
- Follow Go conventions for code structure and formatting
- Write testable code
- Write tests for the code
- Document the architecture and design decisions in a architecture.md file in the docs folder
- Always check if the code still compiles after your changes
- Always run the tests after your changes
- Use test driven development. Write tests first and then write the code to make the tests pass.
- Never use flake utils. Assume the build environment is x86_64-linux.

## Relevant Links
- homeassistant mqtt discovery: https://www.home-assistant.io/integrations/mqtt#mqtt-discovery

## Relevant libraries
- mqtt library: eclipse/paho.mqtt.golang
- filesystem change: fsnotify/fsnotify

## Development environment
- Nix flakes are used to provide a reproducible development environment. You can enter it with `nix develop`.