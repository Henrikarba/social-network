import { socket } from './ws'

export function formatTime(dateStr) {
	const date = new Date(dateStr)
	return date.toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'long',
		day: 'numeric',
		hour: 'numeric',
		minute: 'numeric',
	})
}

export function formatDateTime(dateStr) {
	const date = new Date(dateStr)
	return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })
}

export async function getProfile(id) {
	return new Promise((resolve, reject) => {
		const action = 'get_user'
		const data = {
			user: { id: id },
		}

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
			data: data,
		}

		socket.send(JSON.stringify(requestData))
	})
}

export function isValidFileType(fileType) {
	return fileType === 'image/png' || fileType === 'image/jpeg' || fileType === 'image/jpg' || fileType === 'image/gif'
}
