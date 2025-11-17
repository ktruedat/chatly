import { create } from "zustand";
import { persist } from "zustand/middleware";
import { makeId } from "@/lib/id";

export type ChatMessage = {
  id: string;
  role: "user" | "bot" | "system";
  text?: string;
  imageSrc?: string;
  createdAt: number;
};

export type Chat = {
  id: string;
  title: string;
  createdAt: number;
  updatedAt: number;
  messages: ChatMessage[];
};

type ChatStore = {
  chats: Chat[];
  activeChatId: string | null;
  createChat: (title?: string) => string;
  deleteChat: (id: string) => void;
  selectChat: (id: string) => void;
  addMessage: (chatId: string, msg: ChatMessage) => void;
  upsertMessage: (chatId: string, msg: ChatMessage) => void;
  clearChat: (chatId: string) => void;
  renameChat: (chatId: string, title: string) => void;
};

export const useChatStore = create<ChatStore>()(
  persist<ChatStore>(
    (set: any, get: any) => ({
      chats: [],
      activeChatId: null,
      createChat: (title = "New chat") => {
        const id = makeId();
        const now = Date.now();
        const chat: Chat = {
          id,
          title,
          createdAt: now,
          updatedAt: now,
          messages: [],
        };
        set((s: ChatStore) => ({
          chats: [chat, ...s.chats],
          activeChatId: id,
        }));
        return id;
      },
      deleteChat: (id: string) => {
        set((s: ChatStore) => {
          const chats = s.chats.filter((c) => c.id !== id);
          const activeChatId =
            s.activeChatId === id ? chats[0]?.id ?? null : s.activeChatId;
          return { chats, activeChatId } as Partial<ChatStore> as any;
        });
      },
      selectChat: (id: string) =>
        set(() => ({ activeChatId: id } as Partial<ChatStore> as any)),
      addMessage: (chatId: string, msg: ChatMessage) => {
        set(
          (s: ChatStore) =>
            ({
              chats: s.chats.map((c) => {
                if (c.id !== chatId) return c;
                if (c.messages.some((m) => m.id === msg.id)) return c;
                const isNearDuplicate = c.messages.some(
                  (m) =>
                    m.role === msg.role &&
                    (m.text ?? null) === (msg.text ?? null) &&
                    (m.imageSrc ?? null) === (msg.imageSrc ?? null)
                );
                if (isNearDuplicate) return c;
                return {
                  ...c,
                  messages: [...c.messages, msg],
                  updatedAt: Date.now(),
                };
              }),
            } as Partial<ChatStore> as any)
        );
      },
      upsertMessage: (chatId: string, msg: ChatMessage) => {
        set(
          (s: ChatStore) =>
            ({
              chats: s.chats.map((c) => {
                if (c.id !== chatId) return c;
                const idx = c.messages.findIndex((m) => m.id === msg.id);
                if (idx === -1) {
                  if (msg.id.endsWith("-final")) {
                    const reqId = msg.id.slice(0, -"-final".length);
                    const partialId = reqId + "-partial";
                    const pIdx = c.messages.findIndex(
                      (m) => m.id === partialId
                    );
                    if (pIdx !== -1) {
                      const copy = [...c.messages];
                      copy[pIdx] = { ...copy[pIdx], ...msg };
                      copy[pIdx].id = msg.id;
                      return { ...c, messages: copy, updatedAt: Date.now() };
                    }
                    return {
                      ...c,
                      messages: [...c.messages, msg],
                      updatedAt: Date.now(),
                    };
                  }

                  const isNearDuplicate = c.messages.some(
                    (m) =>
                      m.role === msg.role &&
                      (m.text ?? null) === (msg.text ?? null) &&
                      (m.imageSrc ?? null) === (msg.imageSrc ?? null)
                  );
                  if (isNearDuplicate) return c;
                  return {
                    ...c,
                    messages: [...c.messages, msg],
                    updatedAt: Date.now(),
                  };
                }

                const copy = [...c.messages];
                copy[idx] = { ...copy[idx], ...msg };
                return { ...c, messages: copy, updatedAt: Date.now() };
              }),
            } as Partial<ChatStore> as any)
        );
      },
      clearChat: (chatId: string) => {
        set(
          (s: ChatStore) =>
            ({
              chats: s.chats.map((c) =>
                c.id === chatId
                  ? { ...c, messages: [], updatedAt: Date.now() }
                  : c
              ),
            } as Partial<ChatStore> as any)
        );
      },
      renameChat: (chatId: string, title: string) => {
        set(
          (s: ChatStore) =>
            ({
              chats: s.chats.map((c) =>
                c.id === chatId ? { ...c, title, updatedAt: Date.now() } : c
              ),
            } as Partial<ChatStore> as any)
        );
      },
    }),
    {
      name: "chatly.chats.v1",

      partialize: (state: ChatStore) =>
        ({ chats: state.chats, activeChatId: state.activeChatId } as any),
    }
  )
);

export default useChatStore;
