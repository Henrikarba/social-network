import { writable } from 'svelte/store'

export const currentChat = writable(null)
export const messagesStore = writable(null)
export const groupMessagesStore = writable(null)
