<script>
	import { fly, fade } from 'svelte/transition'
	import { socket, createWebSocket } from '../ws'
	const xplogo = new URL('../assets/xplogo2.png', import.meta.url).href

	let email
	let password
	let logging = false

	async function login() {
		logging = true
		console.log(email, password)
		const creds = {
			email: email,
			password: password,
		}
		try {
			const response = await fetch('http://localhost:80/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(creds),
				credentials: 'include', // Important
			})
			if (!response.ok) {
				throw new Error('Network response was not ok')
			}
			createWebSocket()
		} catch (error) {
			console.error('Error:', error)
		}
	}
</script>

{#if logging}
	<div class="alert bg-red-500">
		<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-info shrink-0 w-6 h-6"
			><path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
			/></svg
		>
		<span>logging innnnnn</span>
	</div>
{/if}
<main>
	<div class=" h-64 w-4/12 select-none" in:fade|global={{ duration: 500 }} out:fly|global={{ duration: 300, x: 800 }}>
		<div
			class="h-10 bg-gradient-to-t from-blue-500 to-blue-700 flex justify-center items-center border-2 rounded border-blue-950"
		>
			<h2 class="text-white font-extrabold">Log on to Windows</h2>
		</div>
		<div class="h-1/3 bg-gray-500 logo flex justify-center border-b-2 border-stone-700">
			<img src={xplogo} class="h-20" alt="XP LOGO" />
		</div>
		<div
			class="h-2/3 bg-gray-500 flex flex-col items-center justify-center gap-2 border-b-2 border-gray-500 rounded-b-lg"
		>
			<div class="flex items-center gap-2">
				<label for="email" class="w-24 text-right">Email</label>
				<input
					bind:value={email}
					type="email"
					placeholder="Type here"
					class="input w-full max-w-xs h-8 bg-white focus:outline-none"
				/>
			</div>
			<div class="flex items-center gap-2">
				<label for="password" class="w-24 text-right">Password</label>
				<input
					bind:value={password}
					type="password"
					placeholder="Type here"
					class="input w-full max-w-xs h-8 bg-white focus:outline-none"
				/>
			</div>
			<div>
				<button class="btn" on:click={login}>LOGIN</button>
			</div>
		</div>
	</div>
</main>

<style lang="scss">
	main {
		background-color: #004e98;
		z-index: 99999999999999;
		height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;

		.logo {
			background: linear-gradient(to right, #6286e1, #6286e1, #9fbbf6, #7899e9);
		}
	}
</style>
