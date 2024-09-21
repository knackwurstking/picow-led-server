import { CleanUp, html, UIStackLayoutPage } from "ui";
import { DialogDeviceSetup } from "../../components";
import { devicesEvents } from "../../lib";
import { DeviceItem } from "./device-item";

export class DevicesPage extends UIStackLayoutPage {
    static register = () => {
        DeviceItem.register();
        DialogDeviceSetup.register();

        customElements.define("devices-page", DevicesPage);
    };

    constructor() {
        super("devices");

        this.cleanup = new CleanUp();

        /** @type {PicowStore} */
        this.uiStore = document.querySelector("ui-store");

        /** @type {import("../../components/picow-app-bar").PicowAppBar} */
        this.picowAppBar = document.querySelector("picow-app-bar");

        this.render();
    }

    render() {
        this.classList.add("no-scrollbar");
        this.shadowRoot.innerHTML += html`
            <style>
                :host {
                    padding-top: var(--ui-app-bar-height);
                    overflow: auto;
                }
            </style>
        `;

        this.innerHTML = html` <ul></ul> `;

        this.picowAppBar.picow.events.on("add", () => {
            const dialog = new DialogDeviceSetup();

            dialog.ui.events.on("close", () => {
                document.body.removeChild(dialog);
            });

            dialog.ui.events.on("submit", async (device) => {
                const s = this.uiStore.ui.get("server");
                const addr = !s.port ? s.host : `${s.host}:${s.port}`;
                const url = `${s.ssl ? "https:" : "http:"}//${addr}/api/device`;
                const r = await fetch(url, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(device),
                });

                if (!r.ok) {
                    // TODO: Add and "error" alert
                    r.text().then((r) => console.error(r));
                    console.error(
                        `Fetch from "${url}" with status code ${r.status}`
                    );
                }
            });

            document.body.appendChild(dialog);
            dialog.ui.open(true);
        });
    }

    connectedCallback() {
        super.connectedCallback();

        this.cleanup.add(
            this.uiStore.ui.on(
                "devices",
                (devices) => {
                    const ul = this.querySelector("ul");

                    while (!!ul.firstChild) ul.removeChild(ul.firstChild);

                    for (const d of devices) {
                        ul.appendChild(new DeviceItem(d));
                    }
                },
                true
            ),

            devicesEvents.events.on("open", () => {
                this.fetchDevices();
            }),

            devicesEvents.events.on("message", (devices) => {
                this.uiStore.ui.set("devices", devices);
            })
        );

        this.fetchDevices();
    }

    disconnectedCallback() {
        super.disconnectedCallback();
        this.cleanup.run();
    }

    async fetchDevices() {
        const s = this.uiStore.ui.get("server");
        const addr = !s.port ? s.host : `${s.host}:${s.port}`;
        const url = `${s.ssl ? "https:" : "http:"}//${addr}/api/devices`;
        const r = await fetch(url, {
            method: "GET",
        });
        if (!r.ok) {
            // TODO: Add and "error" alert
            r.text().then((r) => console.error(r));
            console.error(`Fetch from "${url}" with status code ${r.status}`);
            return;
        }

        this.uiStore.ui.set("devices", await r.json());
    }
}
