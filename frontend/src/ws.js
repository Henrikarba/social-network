import { postsStore, groupPostsStore } from './stores/post'

export const socket = new WebSocket('ws://localhost:80/ws')

socket.onmessage = (msg) => {
	const newData = JSON.parse(msg.data)

	postsStore.update(($postsStore) => ($postsStore = newData.feed.posts))
	groupPostsStore.update(($groupPostsStore) => ($groupPostsStore = newData.feed.group_posts))
}

socket.onopen = () => {
	console.log('WebSocket Connection established')
}
socket.onerror = (error) => {
	console.error('WebSocket connection error:', error)
}
