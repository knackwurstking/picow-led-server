import * as deviceItem from "./device-item";

document.addEventListener("DOMContentLoaded", async () => {
    setupAppBar();

    window.store.listen("devices", (data) => {
        const devices = data;
        const devicesList = document.querySelector<HTMLElement>(
            "._content.devices > .list",
        )!;
        devicesList.innerHTML = "";

        devices.forEach((device) => {
            const item = deviceItem.create(device, () => {
                powerButtonToggle(device);
            });

            devicesList.appendChild(item);
        });
    });

    await window.api.devices();
    setTimeout(() => {
        window.ws.events.addListener("open", async () => {
            await window.api.devices();
        });
    });
});

async function setupAppBar() {
    const items = window.utils.setupAppBarItems(
        "online-indicator",
        "title",
        "settings-button",
    );

    items["title"]!.innerText = "Devices";
}

async function powerButtonToggle(device: Device) {
    let color: Color;
    if (Math.max(...(device.color || [])) > 0) {
        color = (device.pins || device.color || []).map(() => 0);
    } else {
        color = currentColorForDevice(device);
    }

    await window.api.setDevicesColor(color, device);
}

function currentColorForDevice(device: Device): Color {
    return (
        window.store.currentDeviceColor(device.server.addr) ||
        (device.pins || []).map(() => 255)
    );
}
