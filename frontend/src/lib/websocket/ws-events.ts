import { Events } from "ui";
import { BaseWebSocketEvents } from "./base-web-socket-events";

import * as types from "@types";

export class WSEvents extends BaseWebSocketEvents {
    events: Events<{
        server: types.WSEventsServer | null;
        open: null;
        close: null;
        message: any;
        "message-devices": types.WSEventsDevice[];
        "message-error": string;
        "message-device": types.WSEventsDevice;
    }> = new Events();

    constructor() {
        super("/ws");
    }

    get server() {
        return super.server;
    }

    set server(value) {
        super.server = value;
        this.events.dispatch("server", value);
    }

    async request<T extends keyof types.WSEventsCommand>(
        command: T,
        data?: types.WSEventsCommand[T],
    ) {
        if (!this.isOpen()) return;
        console.debug(`[ws] Send command: "${command}"`, {
            server: this.server,
            data,
        });

        let request: types.WSEventsRequest = {
            command: command,
            data: data === undefined ? undefined : JSON.stringify(data),
        };

        switch (command) {
            case "GET api.devices":
            case "POST api.device":
            case "PUT api.device":
            case "DELETE api.device":
            case "POST api.device.pins":
            case "POST api.device.color":
                this.ws?.send(JSON.stringify(request));
                break;

            default:
                throw new Error(`unknown command ${command}`);
        }
    }

    async handleMessageEvent(ev: MessageEvent) {
        super.handleMessageEvent(ev);
        console.debug("[ws] message.event:", ev);

        if (typeof ev.data === "string") {
            try {
                const resp = JSON.parse(ev.data) as types.WSEventsResponse;
                console.debug(`[ws] message:`, resp);

                switch (resp.type) {
                    case "devices":
                    case "device":
                        this.events.dispatch(`message-${resp.type}`, resp.data);
                        break;
                    case "error":
                        this.events.dispatch(`message-error`, resp.data);
                        break;
                }
            } catch (err) {
                console.warn("[ws] Parsing JSON:", err);
            }
        }

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
