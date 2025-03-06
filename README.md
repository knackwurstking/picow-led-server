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
- [Notes](#notes)
- [TS Types](#ts-types)

<!--toc:end-->

## Getting Started

**TODO**: Getting Started section here

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
