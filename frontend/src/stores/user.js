import { writable } from 'svelte/store'

export const currentUser = writable(null)
export const currentUserGroups = writable(null)
export const currentUserFollowers = writable(null)
export const currentUserFollowing = writable(null)
export const notificationsStore = writable(null)
