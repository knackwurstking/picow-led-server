import * as store from "../../../lib/store";
import * as ws from "../../../lib/ws";

export const color = {
    set(addr: string, color?: number[] | null): void {
        if (!color) {
            return;
        }

        if (!color.find((c) => c > 0)) {
            return;
        }

        store.obj.update("color", (data) => {
            data.devices[addr] = color;
            return data;
        });
    },

    get(device: ws.WSDevice) {
        let color = store.obj.get("color")?.devices[device.server.addr];
        if (!color) {
            color = device.pins?.map(() => 255) || [255, 255, 255];
        }
        return color;
    },
};

export function getPowerButtonColor(color?: number[] | null) {
    if (!color) {
        return `rgb(0, 0, 0)`;
    }

    const diff = 255 - Math.max(...color);

    color = color.map((c) => {
        if (c === 0) {
            return c;
        }

        return c + diff;
    });

    return `rgb(${color[0] || 0} ,${color[1] || 0} ,${color[2] || 0})`;
}
