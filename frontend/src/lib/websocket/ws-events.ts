import { Events } from "ui";
import { BaseWebSocketEvents } from "./base-web-socket-events";

type WSEvents_T = {
    "api.devices": Device[];
};

export class WSEvents extends BaseWebSocketEvents {
    events: Events<{
        server: Server | null;
        open: null;
        close: null;
        message: any;
        messageDevice: Device;
        messageDevices: Device[];
    }>;

    constructor() {
        super("/ws");
        this.events = new Events();
    }

    get server() {
        return super.server;
    }

    set server(value) {
        super.server = value;
        this.events.dispatch("server", value);
    }

    async get<T extends keyof WSEvents_T>(path: T): Promise<WSEvents_T[T]> {
        switch (path) {
            case "api.devices":
                // TODO: Send a "GET api.devices" to the server
                break;
        }

        throw new Error(`unknown path ${path}`);
    }

    async handleMessageEvent(ev: MessageEvent) {
        super.handleMessageEvent(ev);
        console.debug("[ws] event:", ev);
        console.debug("[ws] data:", ev.data);
        // TODO: Parsing data and dispatch "message-device" or "message-devices"

        //if (ev.data instanceof Blob) {
        //    this.ws.send("pong");
        //    return;
        //}

        //const device = JSON.parse(ev.data);
        //this.events.dispatch("message", device);
        this.events.dispatch("message", ev.data);
    }

    async handleOpenEvent() {
        await super.handleOpenEvent();
        this.events.dispatch("open", null);
    }

    async handleCloseEvent() {
        await super.handleCloseEvent();
        this.events.dispatch("close", null);
    }
}
