import { writable } from 'svelte/store'
import { currentUser, currentUserGroups, currentUserFollowers, currentUserFollowing } from './stores/user'
import { postsStore, groupPostsStore } from './stores/post'

export let socket

export const isAuthenticated = writable(false)

export function createWebSocket() {
	const ws = new WebSocket('ws://localhost:80/ws')
	ws.onmessage = (msg) => {
		const newData = JSON.parse(msg.data)
		console.log(newData)
		if (newData?.user) {
			postsStore.update(($postsStore) => ($postsStore = newData.feed.posts))
			groupPostsStore.update(($groupPostsStore) => ($groupPostsStore = newData.feed.group_posts))
			currentUser.update(($currentUser) => ($currentUser = newData.user))
			currentUserGroups.update(($currentUserGroups) => ($currentUserGroups = newData.groups))
			currentUserFollowers.update(($currentUserFollowers) => ($currentUserFollowers = newData.followers))
			currentUserFollowing.update(($currentUserFollowing) => ($currentUserFollowing = newData.following))
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
