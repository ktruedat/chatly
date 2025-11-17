"use client";

import React from "react";
import { Button } from "@/shared/components/ui/button";
import { FiTrash } from "react-icons/fi";
import useChatStore from "@/shared/state/chatStore";

export default function ChatSidebar() {
  const chats = useChatStore((s) => s.chats);
  const activeChatId = useChatStore((s) => s.activeChatId);
  const createNewChat = useChatStore((s) => s.createChat);
  const selectChat = useChatStore((s) => s.selectChat);
  const deleteChat = useChatStore((s) => s.deleteChat);

  return (
    <aside className="flex w-64 flex-col gap-3 border-r border-zinc-800 bg-white dark:bg-[#0b0b0b] p-4 rounded-r-2xl">
      <div>
        <Button
          variant="ghost"
          className="w-full justify-start"
          onClick={() => createNewChat("New chat")}
        >
          + New chat
        </Button>
      </div>
      <nav className="flex-1 overflow-y-auto py-2">
        <ul className="flex flex-col gap-2">
          {chats.length === 0 && (
            <li className="rounded-md px-3 py-2 text-sm text-zinc-400">
              No chats
            </li>
          )}
          {chats.map((c) => (
            <li
              key={c.id}
              onClick={() => selectChat(c.id)}
              className={`group cursor-pointer rounded-md px-3 py-2 text-sm hover:bg-zinc-900 flex items-center justify-between ${
                activeChatId === c.id ? "bg-zinc-900" : "text-zinc-300"
              }`}
            >
              <div className="truncate pr-2">
                <div className="font-medium text-sm">{c.title}</div>
                <div className="text-xs text-zinc-500 truncate">
                  {c.messages.length > 0
                    ? c.messages[c.messages.length - 1].text ?? "[image]"
                    : "No messages"}
                </div>
              </div>
              <div className="flex items-center gap-2">
                <div className="text-xs text-zinc-500 ml-2 hidden group-hover:block">
                  {new Date(c.updatedAt).toLocaleTimeString()}
                </div>
                <button
                  aria-label="Delete chat"
                  onClick={(e) => {
                    e.stopPropagation();
                    deleteChat(c.id);
                  }}
                  className="ml-2 inline-flex h-7 w-7 items-center justify-center rounded hover:bg-zinc-800/60 text-zinc-400"
                >
                  <FiTrash size={14} />
                </button>
              </div>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}
