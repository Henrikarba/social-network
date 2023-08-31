import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vitejs.dev/config/
export default defineConfig({
	server: {
		host: true,
		port: 5000,
	},
	plugins: [svelte()],
})
