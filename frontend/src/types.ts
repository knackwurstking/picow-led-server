import type { UIStore, UIThemeHandlerThemes } from "ui";

import type { WSEventsDevice, WSEventsServer } from "./lib/websocket";

export type PicowStackLayoutPages = "devices" | "settings" | "";

export type PicowStore = UIStore<PicowStoreEvents>;

export interface PicowStoreEvents {
    currentPage: PicowStackLayoutPages;
    devices: WSEventsDevice[];
    devicesColor: { [key: string]: number[] };
    server: WSEventsServer;
    currentTheme: {
        theme: UIThemeHandlerThemes;
    };
}

export interface AppBarEvents {
    menu: Event;
    add: Event;
}
