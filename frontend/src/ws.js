import { currentUser, currentUserGroups, currentUserFollowers, currentUserFollowing } from './stores/user'
import { postsStore, groupPostsStore } from './stores/post'

export const socket = new WebSocket('ws://localhost:80/ws')

socket.onmessage = (msg) => {
	const newData = JSON.parse(msg.data)

	// Initial data
	if (newData?.user) {
		postsStore.update(($postsStore) => ($postsStore = newData.feed.posts))
		groupPostsStore.update(($groupPostsStore) => ($groupPostsStore = newData.feed.group_posts))
		currentUser.update(($currentUser) => ($currentUser = newData.user))
		currentUserGroups.update(($currentUserGroups) => ($currentUserGroups = newData.groups))
		currentUserFollowers.update(($currentUserFollowers) => ($currentUserFollowers = newData.followers))
		currentUserFollowing.update(($currentUserFollowers) => ($currentUserFollowers = newData.following))
	}
}

socket.onopen = () => {
	console.log('WebSocket Connection established')
}
socket.onerror = (error) => {
	console.error('WebSocket connection error:', error)
}
