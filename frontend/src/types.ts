import type { UIStackLayout, UIStore } from "ui";
import type { WSEvents_Device, WSServer } from "./lib/websocket";

export type PicowStackLayout = UIStackLayout<PicowStackLayout_Pages>;
export type PicowStackLayout_Pages = null | "devices" | "settings";

export type PicowStore = UIStore<PicowStore_Events>;
export interface PicowStore_Events {
    devices: WSEvents_Device[];
    currentPage: PicowStackLayout_Pages;
    server: WSServer;
}
