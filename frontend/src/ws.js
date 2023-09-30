import { writable } from 'svelte/store'
import {
	currentUser,
	currentUserGroups,
	notificationsStore,
	currentUserFollowers,
	currentUserFollowing,
	chatStore,
	eventStore,
} from './stores/user'
import { currentChat, groupMessagesStore, messagesStore } from './stores/chat'
import { postsStore, groupPostsStore } from './stores/post'

export let socket

export const isAuthenticated = writable(false)

export function createWebSocket() {
	const ws = new WebSocket('ws://localhost:80/ws')
	ws.onmessage = (msg) => {
		const newData = JSON.parse(msg.data)
		if (newData?.user) {
			postsStore.update(($postsStore) => ($postsStore = newData.feed.posts))
			groupPostsStore.update(($groupPostsStore) => ($groupPostsStore = newData.feed.group_posts))
			currentUser.update(($currentUser) => ($currentUser = newData.user))
			currentUserGroups.update(($currentUserGroups) => ($currentUserGroups = newData.groups))
			currentUserFollowers.update(($currentUserFollowers) => ($currentUserFollowers = newData.followers))
			currentUserFollowing.update(($currentUserFollowing) => ($currentUserFollowing = newData.following))
			notificationsStore.update(($notificationsStore) => ($notificationsStore = newData.notifications))
			chatStore.update(($chatStore) => ($chatStore = newData.chatlist))
			messagesStore.update(($messagesStore) => ($messagesStore = newData.messages))
			groupMessagesStore.update(($groupMessagesStore) => ($groupMessagesStore = newData.group_messages))
			eventStore.update(($eventStore) => ($eventStore = newData.events))
		} else if (newData.action == 'get_chat' || newData.action == 'get_group_chat') {
			currentChat.update(($currentChat) => ($currentChat = newData.data))
		} else if (newData.type == 'group') {
			messagesStore.update(($messagesStore) => {
				if (!$messagesStore) $messagesStore = []
				return [...$messagesStore, newData]
			})
		}
	}

	ws.onopen = () => {
		isAuthenticated.set(true)
		socket = ws
		console.log('WebSocket Connection established')
	}
	ws.onerror = (error) => {
		isAuthenticated.set(false)
		console.error('WebSocket connection error:', error)
	}
	ws.onclose = () => {
		isAuthenticated.set(false)
	}
}
