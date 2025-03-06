# PicoW LED Server

Here's the updated project description with the correct name:

**Project Description**:

Picow LED Server is a web-based control system that manages and coordinates multiple Raspberry
Pi Pico W microcontrollers to control LED lighting installations.
The system provides a centralized web interface where users can control LED parameters such
as brightness, color, and patterns across different connected Pico W devices.

**Key Features**:

- Web interface for real-time LED control and monitoring
- Multiple Pico W device management and coordination
- Custom lighting pattern programming and scheduling
- RESTful API for device communication

> Ok, Now for real, this description was written using AI.
> Without this project, I have no lights at home, but until v1.0.0 is released, it is still work in
> progress but will always work somehow.

**Projects belonging to this repository**:

- [PicoW LED Microcontroller](https://github.com/knackwurstking/picow-led-microcontroller)

# Table of Contents

<!--toc:start-->

- [Getting Started](#getting-started)
    - [System Requirements](#system-requirements)
    - [Installation](#installation)
    - [Configuration](#configuration)
    - [Usage](#usage)
- [Notes](#notes)
    - [TS Types](#ts-types)

<!--toc:end-->

## Getting Started

### System Requirements

- Node.js v18+ (Not sure about this, I'm on version v23.9.0)
- Golang v1.24+

### Installation

1. Clone this repository:

    ```bash
    git clone https://github.com/knackwurstking/picow-led-server.git
    cd picow-led-server
    ```

2. Install dependencies:

    ```bash
    make init
    ```

3. Have a look in the makefile and edit the first line to match your configuration location

    ```makefile
    CONFIG_LOCATION=.api.dev.json
    ```

4. Edit the `picow-led-server.service` file to match your configuration location if you want
   to run it as a systemd user service

    ```
    [Service]
    ExecStart=picow-led-server -d -c %h/.config/picow-led-server/api.json
    ```

5. Build the executable with, the executable will be located int the `dist` folder

    ```bash
    make build
    ```

6. Or start the backend development server

    ```bash
    make dev
    ```

7. And start the frontend development server

    ```bash
    cd ui && npm run dev
    ```

8. See also: [picow-led-microcontroller](https://github.com/knackwurstking/picow-led-microcontroller)

### Configuration

Create a `config.json` file in the root directory or wherever you want:

```json
{
    "devices": [
        {
            "server": {
                "name": "Kitchen",
                "addr": "192.168.178.58:3000",
                "online": true
            },
            "pins": [0, 1, 2, 3]
        },
        {
            "server": {
                "name": "Living Room",
                "addr": "192.168.178.50:3000",
                "online": true
            },
            "pins": [0, 1, 2, 3]
        },
        {
            "server": {
                "name": "PC Room",
                "addr": "192.168.178.68:3000",
                "online": true
            },
            "pins": [0, 1, 2, 3]
        },
        {
            "server": {
                "name": "Bath Room",
                "addr": "192.168.178.62:3000",
                "online": true
            },
            "pins": [0, 1, 2, 3]
        },

        {
            "server": {
                "name": "Work Room",
                "addr": "192.168.178.54:3000",
                "online": true
            },
            "pins": [0, 1, 3, 3]
        },
        {
            "server": {
                "name": "Sleep Room",
                "addr": "192.168.178.67:3000",
                "online": true
            },
            "pins": [0, 1, 2, 3]
        }
    ]
}
```

Add this to the command line `-c <config-path>`

### Usage

1. Access the web interface at `http://localhost:50833` or if frontend dev server is running `http://localhost:5173`

## Notes

**Endpoints**:

| Endpoint | GET | POST | PUT | DELETE |
| -------- | :-: | :--: | :-: | :----: |
| /        |  x  |      |     |        |
| /ws      |  x  |      |     |        |

**WebSocket Events**:

| Type    | Data                      |
| ------- | ------------------------- |
| error   | `string`                  |
| devices | [`WSDevice[]`](#ts-types) |
| device  | [`WSDevice`](#ts-types)   |

**WebSocket Commands**:

| Command               | Data                                |
| --------------------- | ----------------------------------- |
| GET api.devices       | `null`                              |
| POST api.device       | [`WSDevice`](#ts-types)             |
| PUT api.device        | [`WSDevice`](#ts-types)             |
| DELETE api.device     | `{ addr: string }`                  |
| POST api.device.pins  | `{ addr: string; pins: number[] }`  |
| POST api.device.color | `{ addr: string; color: number[] }` |

### TS Types

_[ui/src/lib/ws/types.ts](ui/src/lib/ws/types.ts)_

```typescript
export interface WSDeviceServer {
    name?: string;
    addr: string;
    online?: boolean;
}

export interface WSDevice {
    server: WSDeviceServer;
    pins?: number[];
    color?: number[];
}

export type WSCommand = {
    "GET api.devices": null;
    "POST api.device": WSDevice;
    "PUT api.device": WSDevice;
    "DELETE api.device": { addr: string };
    "POST api.device.pins": { addr: string; pins: number[] };
    "POST api.device.color": { addr: string; color: number[] };
};

export interface WSRequest {
    command: string;
    data?: string; // JSON string
}

export type WSResponse =
    | {
          data: string;
          type: "error";
      }
    | {
          data: WSDevice[];
          type: "devices";
      }
    | {
          data: WSDevice;
          type: "device";
      };
```
