# PicoW LED Server

**TODO**: Short project description here

<!--toc:start-->

- [PicoW LED Server](#picow-led-server)
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
