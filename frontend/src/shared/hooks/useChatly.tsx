"use client";

import { useEffect, useRef, useState } from "react";
import ChatlyClient from "@/shared/services/chatlyClient";
import type {
  ClientMessage,
  ServerMessage,
  SubmitRequest,
  CancelRequest,
  ChatMessage,
} from "@/shared/types/ws";
import { makeId } from "@/lib/id";

export function useChatly(initialAutoConnect = true) {
  const clientRef = useRef<ChatlyClient | null>(null);
  const [connected, setConnected] = useState(false);
  const [serverMessages, setServerMessages] = useState<ServerMessage[]>([]);
  const [chatMessages, setChatMessages] = useState<ChatMessage[]>([]);

  useEffect(() => {
    clientRef.current = new ChatlyClient();

    const onOpen = () => setConnected(true);
    const onClose = () => setConnected(false);
    const onError = () => setConnected(false);

    const onMessage = (m: ServerMessage) => {
      setServerMessages((s) => [...s, m]);

      if (m.type === "final_response") {
        setChatMessages((prev) => {
          const partialId = m.request_id + "-partial";
          const idx = prev.findIndex((p) => p.id === partialId);
          const finalMsg = {
            id: m.request_id + "-final",
            role: "bot",
            text: m.response_text,
            partial: false,
          } as ChatMessage;
          if (idx !== -1) {
            const copy = [...prev];
            copy[idx] = finalMsg;
            return copy;
          }
          return [...prev, finalMsg];
        });
      } else if (m.type === "token") {
        setChatMessages((prev) => {
          const partialId = m.request_id + "-partial";
          const idx = prev.findIndex((p) => p.id === partialId);
          if (idx !== -1) {
            const updated = {
              ...prev[idx],
              text: (prev[idx].text || "") + m.token,
            } as ChatMessage;
            const copy = [...prev];
            copy[idx] = updated;
            return copy;
          }
          return [
            ...prev,
            {
              id: partialId,
              role: "bot",
              text: m.token,
              partial: true,
            } as ChatMessage,
          ];
        });
      } else if (m.type === "ack") {
        setChatMessages((prev) => [
          ...prev,
          {
            id: makeId(),
            role: "system",
            text: `Ack: ${m.status}`,
          } as ChatMessage,
        ]);
      } else if (m.type === "error") {
        setChatMessages((prev) => [
          ...prev,
          {
            id: makeId(),
            role: "system",
            text: `Error: ${m.message}`,
          } as ChatMessage,
        ]);
      } else if (m.type === "progress") {
        setChatMessages((prev) => [
          ...prev,
          { id: makeId(), role: "system", text: m.message } as ChatMessage,
        ]);
      }
    };

    clientRef.current.onOpen(onOpen);
    clientRef.current.onClose(onClose);
    clientRef.current.onError(onError);
    clientRef.current.onMessage(onMessage);

    if (initialAutoConnect) {
      clientRef.current
        .connect()
        .then(() => setConnected(true))
        .catch(() => setConnected(false));
    }

    return () => {
      clientRef.current?.disconnect();
      clientRef.current = null;
    };
  }, [initialAutoConnect]);

  const connect = async () => {
    if (!clientRef.current) clientRef.current = new ChatlyClient();
    await clientRef.current.connect();
    setConnected(true);
  };

  const disconnect = () => {
    clientRef.current?.disconnect();
    setConnected(false);
  };

  const sendSubmitRequest = async (data: {
    input_type: "Text" | "Image" | "ImageHard";
    content?: string | null;
    image_base64?: string | null;
    metadata?: Record<string, any> | null;
    request_id?: string;
  }) => {
    const request_id = data.request_id ?? makeId();
    const msg: SubmitRequest = {
      type: "submit_request",
      request_id,
      input_type: data.input_type,
      content: data.content ?? null,
      image_base64: data.image_base64 ?? null,
      metadata: data.metadata ?? null,
    };

    if (!clientRef.current) clientRef.current = new ChatlyClient();

    if (!clientRef.current.isConnected()) {
      try {
        await clientRef.current.connect();
        setConnected(true);
      } catch (e) {
        throw e;
      }
    }

    clientRef.current.sendMessage(msg as ClientMessage);

    setChatMessages((prev) => [
      ...prev,
      {
        id: request_id + "-user",
        role: "user",
        text: data.content ?? undefined,
        imageSrc: data.image_base64
          ? `data:image;base64,${data.image_base64}`
          : undefined,
      },
    ]);

    return request_id;
  };

  const sendText = async (text: string) => {
    return sendSubmitRequest({ input_type: "Text", content: text });
  };
  const sendTextWithId = async (text: string, request_id?: string) => {
    return sendSubmitRequest({ input_type: "Text", content: text, request_id });
  };

  const sendImageFile = (
    file: File,
    input_type: "Image" | "ImageHard" = "Image",
    request_id?: string,
    content?: string | null
  ) => {
    return new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = async () => {
        const result = reader.result as string;
        const base64 = result
          .replace(/^data:\w+\/[-+.\w]+;base64,/, "")
          .replace(/^data:[^;]+;base64,/, "");
        try {
          const id = await sendSubmitRequest({
            input_type,
            image_base64: base64,
            content: content ?? null,
            request_id,
          });
          resolve(id);
        } catch (e) {
          reject(e);
        }
      };
      reader.onerror = (e) => reject(e);
      reader.readAsDataURL(file);
    });
  };

  const sendCancel = (request_id: string) => {
    const msg: CancelRequest = { type: "cancel_request", request_id };
    if (!clientRef.current) clientRef.current = new ChatlyClient();
    if (!clientRef.current.isConnected()) {
      clientRef.current
        .connect()
        .then(() => clientRef.current?.sendMessage(msg));
    } else {
      clientRef.current.sendMessage(msg);
    }
  };

  return {
    connected,
    connect,
    disconnect,
    sendText,
    sendTextWithId,
    sendImageFile,
    sendCancel,
    serverMessages,
    chatMessages,
  };
}

export default useChatly;
