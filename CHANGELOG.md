# CHANGELOG

## v0.10.0 [unreleased]

**Changed**:

- Replaced gorilla with [golang.org/x/net/websocket]
- Update ui dependency to v2.1.0

## v0.9.2 [2025-02-25]

**Misc**:

- Rename "frontend" to "ui"
- Code clean up

**Fixed**:

- Device handling while unmarshalling API data, if devices are offline.

## v0.9.1 [2025-02-18]

- Fixed missing color check for null

## v0.9.0 [2025-02-13]

**Updated**:

- Updated UI library to [v2.0.0](https://github.com/knackwurstking/ui/tree/dev?tab=readme-ov-file)

## v0.8.3 [2024-11-18]

**Fixed**:

- Fixed orientation lock on mobile, remove orientation from manifest
- Fixed html closing tag

## v0.8.2 [2024-11-18]

[ui v1.1.0](https://github.com/knackwurstking/ui)

**Removed**:

- Removed scrollbars from pages

## v0.8.1 [2024-10-24 Quick Fix]

**Fixed**:

- Fixed wrong fonts path for `includeAssets` (PWA manifest), and themes added

## v0.8.0 [2024-10-24]

**General**:

- [frontend] Updated
  ["ui" to version v1.0.0](https://github.com/knackwurstking/ui)
  , using Lit web components now
- [backend] Moved "ws" package to models/

## v0.7.0 [2024-10-09]

- [ui v0.8.0](https://github.com/knackwurstking/ui)

**General**:

- [backend] Completely rewritten, using only websockets now
- [frontend] Converted to typescript, updated to fit the new backend

## v0.6.2 [2024-09-30]

**Fixed**:

- [frontend] Fix missing setup-device-dialog content
- [frontend] Fix missing setup-device-dialog action buttons if
  `allowDeletion` was set to false

## v0.6.1 [2024-09-30]

**Added**:

- [backend] Added `events` package ("pkg/event")

**Changed**:

- [frontend] Changed "ui" version in use to v0.7.1

**General**:

- [backend] Moved `endpoints` package to "internal/"
- [backend] Use `events.Events` for API changes for emitting WebSocket events
