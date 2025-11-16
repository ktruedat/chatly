"use client";

import React from "react";
import { Button } from "@/shared/components/ui/button";
import { FiPaperclip, FiSend } from "react-icons/fi";
import useChatView from "@/shared/hooks/useChatView";
import ChatSidebar from "./ChatSidebar";
import { Input } from "@/shared/components/ui/input";

type Message = {
  id: string;
  from: "user" | "bot" | "system";
  type: "text" | "image";
  content?: string;
  imageSrc?: string;
};

export default function ChatView() {
  const {
    connected,
    activeChat,
    text,
    setText,
    previewSrc,
    fileRef,
    handleFileChange,
    handleSend,
    isEmpty,
    typed,
    listRef,
  } = useChatView();

  return (
    <div className="flex h-screen w-full bg-white text-black dark:bg-[#0b0b0b] dark:text-white">
      <ChatSidebar />

      <main className="flex flex-1 flex-col">
        <header className="flex items-center justify-between px-6 py-4">
          <div className="flex items-center gap-3">
            <div className="h-8 w-8 rounded-full bg-white/10" />
            <h1 className="text-lg font-semibold">Chatly</h1>
          </div>
          <div className="flex items-center gap-3">
            <div className="text-sm text-zinc-400">Chatly</div>
            <div className="ml-3 flex items-center gap-2">
              <span
                className={`inline-block h-2 w-2 rounded-full ${
                  connected ? "bg-emerald-400" : "bg-zinc-600"
                }`}
                aria-hidden
              />
              <div className="text-sm text-zinc-400">
                {connected ? "Connected" : "Disconnected"}
              </div>
            </div>
          </div>
        </header>

        <div className="flex-1 overflow-hidden">
          <div ref={listRef} className="h-full overflow-y-auto px-6 py-6">
            <div className="mx-auto w-full max-w-3xl">
              <div className="flex flex-col gap-6">
                {isEmpty ? (
                  <div className="flex h-[60vh] flex-col items-center justify-center gap-6">
                    <div className="rounded-2xl bg-zinc-100/50 p-6 text-center text-zinc-500 dark:bg-zinc-900/40 dark:text-zinc-400">
                      <div className="text-xl font-semibold tracking-widest">
                        <span>{typed}</span>
                        <span className="ml-2 inline-block h-6 w-0.5 bg-zinc-500 dark:bg-zinc-300 animate-pulse align-middle" />
                      </div>
                    </div>
                    {isEmpty && (
                      <div className="w-full max-w-2xl">
                        <div className="rounded-2xl bg-zinc-50 p-4 dark:bg-zinc-900">
                          {previewSrc && (
                            <div className="mb-3 flex items-center gap-3">
                              <div className="relative h-20 w-20 overflow-hidden rounded-2xl bg-zinc-100 dark:bg-zinc-800">
                                <img
                                  src={previewSrc}
                                  alt="preview"
                                  className="h-full w-full object-cover"
                                />
                                <button
                                  aria-label="Remove attached image"
                                  className="absolute right-1 top-1 z-10 inline-flex h-6 w-6 items-center justify-center rounded-full bg-black/60 text-white text-xs hover:bg-black/70"
                                  onClick={() => {
                                    if (previewSrc)
                                      try {
                                        URL.revokeObjectURL(previewSrc);
                                      } catch (e) {}
                                  }}
                                >
                                  ×
                                </button>
                              </div>
                              <div className="flex-1 text-sm text-zinc-300">
                                Image ready to send
                              </div>
                            </div>
                          )}

                          <textarea
                            value={text}
                            onChange={(e) => setText(e.target.value)}
                            onKeyDown={(e) => {
                              if (e.key === "Enter" && !e.shiftKey) {
                                e.preventDefault();
                                handleSend();
                              }
                            }}
                            rows={4}
                            aria-label="Message input"
                            className="w-full resize-none rounded-xl bg-transparent px-3 py-2 text-sm placeholder:text-zinc-500 dark:placeholder:text-zinc-400 focus:outline-none"
                          />
                          <div className="mt-3 flex items-center justify-between">
                            <div className="flex items-center gap-2">
                              <button
                                className="rounded-md p-2 text-zinc-600 hover:bg-zinc-100/50 dark:text-zinc-300 dark:hover:bg-zinc-800/60"
                                onClick={() => fileRef.current?.click()}
                                aria-label="Attach image"
                              >
                                <FiPaperclip size={18} />
                              </button>
                              <input
                                ref={fileRef}
                                onChange={handleFileChange}
                                type="file"
                                accept="image/*"
                                className="hidden"
                              />
                            </div>
                            <div>
                              <button
                                onClick={handleSend}
                                aria-label="Send"
                                className="flex h-10 w-10 items-center justify-center rounded-full bg-[#10a37f] text-white hover:brightness-110"
                              >
                                <FiSend />
                              </button>
                            </div>
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                ) : (
                  [...(activeChat?.messages ?? [])]
                    .filter(
                      (v, i, a) => a.findIndex((x) => x.id === v.id) === i
                    )
                    .map((c) => {
                      const m: Message = {
                        id: c.id,
                        from:
                          c.role === "user"
                            ? "user"
                            : c.role === "bot"
                            ? "bot"
                            : "system",
                        type: c.imageSrc ? "image" : "text",
                        content: c.text ?? undefined,
                        imageSrc: c.imageSrc ?? undefined,
                      };

                      return (
                        <div
                          key={m.id}
                          className={`flex items-start gap-4 ${
                            m.from === "user" ? "justify-end" : "justify-start"
                          }`}
                        >
                          {m.from === "bot" && (
                            <div className="shrink-0">
                              <div className="h-8 w-8 rounded-full bg-zinc-700 flex items-center justify-center text-sm">
                                B
                              </div>
                            </div>
                          )}

                          <div
                            className={`max-w-[80%] ${
                              m.from === "user" ? "text-right" : "text-left"
                            }`}
                          >
                            <div
                              className={`${
                                m.from === "user"
                                  ? "rounded-2xl bg-[#10a37f] text-white"
                                  : "rounded-2xl bg-zinc-900 text-zinc-100"
                              } inline-block px-4 py-3`}
                            >
                              {m.type === "text" && (
                                <div className="whitespace-pre-wrap">
                                  {m.content}
                                </div>
                              )}
                              {m.type === "image" && m.imageSrc && (
                                <div className="flex flex-col items-center gap-2">
                                  <div className="overflow-hidden rounded-2xl bg-black">
                                    <img
                                      src={m.imageSrc}
                                      alt={m.content || "image"}
                                      className="w-full h-auto max-h-96 object-contain"
                                    />
                                  </div>
                                  {m.content && (
                                    <div className="text-sm opacity-80 text-zinc-300">
                                      {m.content}
                                    </div>
                                  )}
                                </div>
                              )}
                            </div>
                          </div>

                          {m.from === "user" && (
                            <div className="shrink-0">
                              <div className="h-8 w-8 rounded-full bg-[#0ea5a1] flex items-center justify-center text-sm">
                                U
                              </div>
                            </div>
                          )}
                        </div>
                      );
                    })
                )}
              </div>
            </div>
          </div>
        </div>

        {!isEmpty && (
          <div className="px-6 py-4">
            <div className="mx-auto flex w-full max-w-3xl items-end gap-3">
              <div className="flex-1">
                {previewSrc && (
                  <div className="mb-3 flex items-center gap-3">
                    <div className="relative h-20 w-20 overflow-hidden rounded-2xl bg-zinc-100 dark:bg-zinc-800">
                      <img
                        src={previewSrc}
                        alt="preview"
                        className="h-full w-full object-cover"
                      />
                      <button
                        aria-label="Remove attached image"
                        className="absolute right-1 top-1 z-10 inline-flex h-6 w-6 items-center justify-center rounded-full bg-black/60 text-white text-xs hover:bg-black/70"
                        onClick={() => {
                          if (previewSrc)
                            try {
                              URL.revokeObjectURL(previewSrc);
                            } catch (e) {}
                        }}
                      >
                        ×
                      </button>
                    </div>
                    <div className="flex-1 text-sm text-zinc-300">
                      Image ready to send
                    </div>
                  </div>
                )}

                <Input
                  value={text}
                  onChange={(e) =>
                    setText((e.target as HTMLInputElement).value)
                  }
                  onKeyDown={(e: React.KeyboardEvent<HTMLInputElement>) => {
                    if (e.key === "Enter" && !e.shiftKey) {
                      e.preventDefault();
                      handleSend();
                    }
                  }}
                  aria-label="Message input"
                  className="min-h-11 w-full rounded-2xl bg-zinc-100 dark:bg-zinc-800/60 px-3 py-2 text-sm placeholder:text-zinc-500 dark:placeholder:text-zinc-400"
                />
              </div>

              <input
                ref={fileRef}
                onChange={handleFileChange}
                type="file"
                accept="image/*"
                className="hidden"
              />

              <div className="flex items-center gap-2">
                <button
                  className="rounded-md p-2 text-zinc-600 hover:bg-zinc-100/50 dark:text-zinc-300 dark:hover:bg-zinc-800/60"
                  onClick={() => fileRef.current?.click()}
                  aria-label="Attach image"
                >
                  <FiPaperclip size={18} />
                </button>
                <button
                  onClick={handleSend}
                  aria-label="Send"
                  className="flex h-10 w-10 items-center justify-center rounded-full bg-[#10a37f] text-white hover:brightness-110"
                >
                  <FiSend />
                </button>
              </div>
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
