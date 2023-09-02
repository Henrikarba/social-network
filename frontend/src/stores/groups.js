import { writable } from 'svelte/store'
import { socket } from '../ws'

export const groupStore = writable(null)

export async function getGroups() {
	return new Promise((resolve, reject) => {
		const action = 'get_groups'

		const temporaryListener = (event) => {
			const response = JSON.parse(event.data)
			if (response.action === action && response.data) {
				socket.removeEventListener('message', temporaryListener)
				resolve(response.data)
			} else if (response.action === 'error') {
				socket.removeEventListener('message', temporaryListener)
				reject(response.error)
			}
		}

		socket.addEventListener('message', temporaryListener)

		const requestData = {
			action: action,
		}

		socket.send(JSON.stringify(requestData))
	})
}

export async function groupInfo(id) {
	return new Promise((resolve, reject) => {
		const action = 'get_group'

		const temporaryListener = (event) => {
			const response = JSON.parse(event.data)
			if (response.action === action && response.data) {
				socket.removeEventListener('message', temporaryListener)
				resolve(response.data)
			} else if (response.action === 'error') {
				socket.removeEventListener('message', temporaryListener)
				reject(response.error)
			}
		}

		socket.addEventListener('message', temporaryListener)

		const requestData = {
			action: action,
			data: {
				id: id,
			},
		}

		socket.send(JSON.stringify(requestData))
	})
}
