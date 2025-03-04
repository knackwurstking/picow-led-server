import "./bootstrap-icons.css";

import { registerSW } from "virtual:pwa-register";

import * as ui from "ui";

import * as constants from "./constants";
import * as store from "./lib/store";
import * as pages from "./pages";

// PWA Updater

const updateSW = registerSW({
    async onNeedRefresh() {
        if (confirm(`Update available, press OK to update.`)) {
            await updateSW();
        }
    },
});

const devicesRoute = {
    title: "PicoW LED | Devices",
    template: {
        selector: `template.devices`,
        onMount: () => pages.devices.onMount(),
        onDestroy: pages.devices.onDestroy,
    },
};

const routes: { [key: string]: ui.router.Route } = {
    "/": devicesRoute,
    devices: devicesRoute,

    settings: {
        title: "PicoW LED | Settings",
        template: {
            selector: `template.settings`,
            onMount: () => pages.settings.onMount(),
            onDestroy: () => pages.settings.onDestroy(),
        },
    },
};

ui.router.hash.init(document.querySelector(`.router-target`)!, routes);

if (store.obj.get("firstTimeConnect")) {
    location.hash = "#settings";
}

document.querySelector<HTMLElement>(`.build`)!.innerText = `${constants.version}`;
