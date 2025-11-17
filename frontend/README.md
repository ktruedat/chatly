# Chatly — Frontend

This folder contains the Chatly frontend application (Next.js). It implements the chat UI, image attach/preview, and a WebSocket client that streams responses from the Chatly backend.

- Framework: Next.js (app router)
- Language: TypeScript
- State: Zustand (persisted to localStorage)
- Styling: Tailwind CSS + shadcn primitives

## Prerequisites

- Node.js or Bun (recommended if you use the bundled `bun.lockb`). The repo was developed with Bun available, but `npm` / `pnpm` also work.
- A running backend that exposes the WebSocket endpoint (default: `http://localhost:5001`). See the backend `README` for backend setup.

## Environment

Create a `.env` file in this folder (a sample is already present). Important env vars:

Mention, you will need to set from .env backend file port to 5001

```

NEXT_PUBLIC_SERVER_URL="http://localhost:5001"
```

- `NEXT_PUBLIC_SERVER_URL` is exposed to the client via Next and should match the backend URL for development.

## Install dependencies

Using Bun (recommended if you have Bun installed):

Move into the frontend folder and then run:

```bash
bun install
```

Or using npm:

```bash
npm install
```

Or using pnpm:

```bash
pnpm install
```

## Run (development)

Start the dev server:

```bash
# from the frontend folder
bun run dev
```

or with npm / pnpm:

```bash
npm run dev
# or
pnpm dev
```

Open http://localhost:3000

## Project architecture

Below is a high-level view of the frontend folder and where the important pieces live. This will help you find the chat UI, the WebSocket client, and the state/store logic quickly.

```
frontend/
├─ public/                # static assets
├─ src/
│  ├─ app/                # Next.js app pages (app router) - page/layout entrypoints
│  │  ├─ layout.tsx
│  │  └─ page.tsx
│  ├─ api/                # REST helpers and interceptors used by services
│  ├─ config/             # configuration constants (e.g. urls)
│  ├─ lib/                # small utilities (id generation, helpers)
│  ├─ shared/
│  │  ├─ components/      # React UI components and pages
│  │  │  └─ pages/         # page-level components (ChatView, ChatSidebar)
│  │  ├─ hooks/           # custom hooks (useChatly, useChatView)
│  │  ├─ services/        # low-level services (ChatlyClient WebSocket wrapper)
│  │  ├─ state/           # Zustand store and persistence (chatStore.ts)
│  │  └─ types/           # shared TypeScript types (WebSocket messages)
│  └─ styles/
└─ package.json
```

### Key files explained

- `src/shared/services/chatlyClient.ts`

  - WebSocket wrapper that opens a connection to the backend, emits/receives messages, and exposes simple `onOpen/onClose/onError/onMessage` hooks. It encapsulates reconnect logic and backoff.

- `src/shared/hooks/useChatly.tsx`

  - Client hook that creates a `ChatlyClient`, registers event handlers, and provides methods to `sendSubmitRequest`, `sendText`, `sendImageFile`, and `sendCancel`. It also exposes `chatMessages` (incoming server messages) and `serverMessages` for debugging.

- `src/shared/hooks/useChatView.tsx`

  - Business-logic hook used by the chat page. It subscribes to `useChatly`, maps request IDs to chat IDs, performs optimistic upserts for user messages, handles file previews, and writes messages into the Zustand `chatStore`.

- `src/shared/state/chatStore.ts`

  - The Zustand store that holds `chats`, `activeChatId`, and message operations (`createChat`, `addMessage`, `upsertMessage`, `deleteChat`, ...). It uses the `persist` middleware to keep chat list and active chat in `localStorage` under the key `chatly.chats.v1`.

- `src/shared/types/ws.ts`
  - TypeScript definitions for the WebSocket protocol used between frontend and backend (token frames, final_response, submit_request, etc.). Use this file as the contract between frontend and backend for message shapes.

### WebSocket protocol (brief)

Client sends: `submit_request` with fields:

- `request_id` (client-generated id to correlate streaming frames)
- `input_type` ("Text" | "Image" | "ImageHard")
- `content` (optional text)
- `image_base64` (optional base64 image)

Server sends streaming frames:

- `token`: partial token for request_id
- `final_response`: final text for request_id
- `ack`, `progress`, `error`

The frontend creates an optimistic user message `requestId-user` on send, collects streaming `token` frames as `requestId-partial` and updates the store with `requestId-final` when final arrives.

## Persistence and localStorage notes

- The app uses Zustand + `persist` so chats survive page reloads. The persisted key is `chatly.chats.v1` in `localStorage`.
