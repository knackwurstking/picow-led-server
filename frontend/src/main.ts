import "../node_modules/ui/css/main.css";

import { html, UIThemeHandler } from "ui";
import { registerSW } from "virtual:pwa-register";

import createAppBar from "./app/create-app-bar";
import createDrawer from "./app/create-drawer";
import PicowDevicesPage from "./app/pages/picow-devices-page";
import PicowSettingsPage from "./app/pages/picow-settings-page";
import { throwAlert } from "./lib/utils";
import ws from "./lib/websocket";
import type {
    PicowStackLayout,
    PicowStackLayout_Pages,
    PicowStore,
} from "./types";

registerSW({
    onRegistered(r) {
        console.debug("ServiceWorkerRegistration:", r);
        if (!r) return;

        setTimeout(async () => {
            try {
                console.debug(`PWA Update service...`);
                await r.update(); // For now do auto update all the time
            } catch (err) {
                console.warn(`PWA Auto update: ${err}`);
            }
        });
    },
});

async function main() {
    const el = document.querySelector<HTMLElement>(`#app`)!;

    el.style.width = "100%";
    el.style.height = "100%";

    el.innerHTML = html`
        <ui-theme-handler mode="dark"></ui-theme-handler>
        <ui-store storageprefix="picow:" storage></ui-store>

        <div class="app-bar"></div>
        <div class="drawer"></div>

        <ui-container style="width: 100%; height: 100%;">
            <ui-stack-layout></ui-stack-layout>
        </ui-container>

        <ui-alerts></ui-alerts>
    `;

    // ---------------- //
    // Create the Store //
    // ---------------- //

    const store = el.querySelector<PicowStore>(`ui-store`);

    {
        el.querySelector<UIThemeHandler>(`ui-theme-handler`).ui.theme =
            store.ui.get("currentTheme")?.theme || null;
    }

    store.ui.set("currentPage", null, true);
    store.ui.set("devices", [], true);

    store.ui.set("devicesColor", {}, true);

    store.ui.set(
        "server",
        {
            ssl: !!location.protocol.match(/(https)/),
            host: location.hostname,
            port: location.port,
        },
        true
    );

    // ---------------------- //
    // Create the StackLayout //
    // ---------------------- //

    const stackLayout = el.querySelector<PicowStackLayout>(`ui-stack-layout`);

    stackLayout.ui.register("devices", async () => {
        return new PicowDevicesPage(appBar);
    });

    stackLayout.ui.register("settings", async () => {
        return new PicowSettingsPage();
    });

    stackLayout.ui.events.on("change", ({ newPage }) => {
        // Reset all layouts (AppBar buttons and title)
        appBar.title = "PicoW LED";
        appBar.items.add.ui.hide();

        if (!newPage) {
            drawer.open();
            return;
        }

        store.ui.set("currentPage", newPage.ui.name as PicowStackLayout_Pages);
        switch (newPage.ui.name) {
            case "devices":
                store.ui.set("currentPage", newPage.ui.name);
                appBar.title = "Devices";
                appBar.items.add.ui.show();
                break;

            case "settings":
                store.ui.set("currentPage", newPage.ui.name);
                appBar.title = "Settings";
                break;

            default:
                appBar.title = newPage.ui.name;
                break;
        }
    });

    // ----------------- //
    // Create the AppBar //
    // ----------------- //

    const appBar = await createAppBar();
    el.querySelector(`div.app-bar`).replaceWith(appBar.element);

    appBar.events.on("menu", () => drawer.open());

    // ----------------- //
    // Create the Drawer //
    // ----------------- //

    const drawer = await createDrawer();
    el.querySelector(`div.drawer`).replaceWith(drawer.element);

    // ---------------- //
    // Lets get started //
    // ---------------- //

    setTimeout(() => {
        // Set the start page
        if (!!store.ui.get("currentPage")) {
            stackLayout.ui.set(store.ui.get("currentPage"));
        } else {
            drawer.open();
        }

        ws.events.on("message-error", (msg) =>
            throwAlert({ message: msg, variant: "error" })
        );

        // Handler server changes
        store.ui.on(
            "server",
            async (server) => {
                ws.server = server;
            },
            true
        );
    });
}

main();
