"use client";

import { useEffect, useRef, useState } from "react";
import useChatly from "@/shared/hooks/useChatly";
import { makeId } from "@/lib/id";
import useChatStore, { ChatMessage } from "@/shared/state/chatStore";

export function useChatView() {
  const { connected, sendTextWithId, sendImageFile, chatMessages } =
    useChatly(true);

  const chats = useChatStore((s) => s.chats);
  const activeChatId = useChatStore((s) => s.activeChatId);
  const createChat = useChatStore((s) => s.createChat);
  const selectChat = useChatStore((s) => s.selectChat);
  const upsertMessage = useChatStore((s) => s.upsertMessage);
  const deleteChat = useChatStore((s) => s.deleteChat);

  const activeChat = chats.find((c) => c.id === activeChatId) ?? null;

  const [text, setText] = useState("");
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [previewSrc, setPreviewSrc] = useState<string | null>(null);
  const WELCOME = "WELCOME TO CHATLY";
  const [typed, setTyped] = useState("");
  const fileRef = useRef<HTMLInputElement | null>(null);
  const listRef = useRef<HTMLDivElement | null>(null);

  const reqToChatRef = useRef<Record<string, string>>({});
  const lastProcessedRef = useRef<number>(0);
  const processedIdsRef = useRef<Set<string>>(new Set());

  const makeIdLocal = () => makeId();

  const handleSend = async () => {
    if (!text.trim() && !imageFile) return;
    let chatId = activeChatId;
    if (!chatId) chatId = createChat("New chat");

    const requestId = makeIdLocal();
    reqToChatRef.current[requestId] = chatId;

    try {
      if (imageFile) {
        const content = text.trim() || null;
        upsertMessage(chatId, {
          id: requestId + "-user",
          role: "user",
          imageSrc: previewSrc ?? undefined,
          text: content ?? undefined,
          createdAt: Date.now(),
        } as ChatMessage);

        if (previewSrc) {
          try {
            URL.revokeObjectURL(previewSrc);
          } catch (e) {
            /* ignore */
          }
        }
        setPreviewSrc(null);
        setImageFile(null);
        if (fileRef.current) fileRef.current.value = "";
        setText("");

        await sendImageFile(imageFile, "Image", requestId, content);
      } else {
        const content = text.trim();
        upsertMessage(chatId, {
          id: requestId + "-user",
          role: "user",
          text: content,
          createdAt: Date.now(),
        } as ChatMessage);

        setText("");
        await sendTextWithId(content, requestId);
      }
    } catch (err) {
      console.error("Failed to send message", err);
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0] || null;
    if (f) {
      if (previewSrc) URL.revokeObjectURL(previewSrc);
      const url = URL.createObjectURL(f);
      setPreviewSrc(url);
      setImageFile(f);
    } else {
      if (previewSrc) URL.revokeObjectURL(previewSrc);
      setPreviewSrc(null);
      setImageFile(null);
    }
  };

  const isEmpty = (activeChat?.messages?.length ?? 0) === 0;

  useEffect(() => {
    if (!chatMessages || chatMessages.length === 0) return;
    for (let i = 0; i < chatMessages.length; i++) {
      const m = chatMessages[i];
      if (processedIdsRef.current.has(m.id)) continue;
      let targetChatId = activeChatId;
      const suffixes = ["-partial", "-final", "-user"];
      const matched = suffixes.find((s) => m.id.endsWith(s));
      if (matched) {
        const reqId = m.id.slice(0, -matched.length);
        const mapped = reqToChatRef.current[reqId];
        if (mapped) targetChatId = mapped;
        if (matched === "-final") {
          try {
            delete reqToChatRef.current[reqId];
          } catch (e) {
            console.error("Failed to delete reqToChat mapping", e);
          }
        }
      }
      if (!targetChatId) continue;
      upsertMessage(targetChatId, {
        id: m.id,
        role: m.role === "user" ? "user" : m.role === "bot" ? "bot" : "system",
        text: m.text ?? undefined,
        imageSrc: m.imageSrc ?? undefined,
        createdAt: Date.now(),
      } as ChatMessage);
      processedIdsRef.current.add(m.id);
    }
  }, [chatMessages]);

  useEffect(() => {
    if (!isEmpty) {
      setTyped("");
      return;
    }
    setTyped("");
    let i = 0;
    let timer: ReturnType<typeof setTimeout> | null = null;
    const tick = () => {
      if (i >= WELCOME.length) return;
      const ch = WELCOME[i];
      if (typeof ch === "string") setTyped((prev) => prev + ch);
      i += 1;
      if (i < WELCOME.length) timer = setTimeout(tick, 80);
    };
    timer = setTimeout(tick, 80);
    return () => {
      if (timer) clearTimeout(timer);
    };
  }, [isEmpty]);

  useEffect(() => {
    return () => {
      if (previewSrc) {
        try {
          URL.revokeObjectURL(previewSrc);
        } catch (e) {
          console.error("Failed to revoke object URL", e);
        }
      }
    };
  }, [previewSrc]);

  const createNewChat = (title = "New chat") => {
    const id = createChat(title);
    setText("");
    if (previewSrc) {
      try {
        URL.revokeObjectURL(previewSrc);
      } catch (e) {
        console.error("Failed to revoke object URL", e);
      }
    }
    setPreviewSrc(null);
    setImageFile(null);
    if (fileRef.current) fileRef.current.value = "";
    if (id) selectChat(id);
    return id;
  };

  return {
    connected,
    chats,
    activeChat,
    activeChatId,
    createNewChat,
    selectChat,
    deleteChat,
    text,
    setText,
    imageFile,
    previewSrc,
    fileRef,
    handleFileChange,
    handleSend,
    isEmpty,
    typed,
    listRef,
  };
}

export default useChatView;
