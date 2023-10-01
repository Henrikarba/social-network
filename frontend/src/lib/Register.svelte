<script>
	import { fly, fade } from 'svelte/transition'
	import { isValidFileType } from '../utils'
	import { socket, createWebSocket } from '../ws'
	import { LoginPage } from '../stores/signin'
	const xplogo = new URL('../assets/xplogo2.png', import.meta.url).href

	let Email
	let Password
	let FirstName
	let LastName
	let DateOfBirth
	let Nickname
	let AboutMe
	let registering = false
	let msg = ''
	let files
	async function register() {
		const formData = new FormData()
		registering = true
		msg = '........registering....'
		formData.append('email', Email)
		formData.append('password', Password)
		formData.append('first_name', FirstName)
		formData.append('last_name', LastName)
		formData.append('about_me', AboutMe)
		formData.append('nickname', Nickname)
		formData.append('date_of_birth', DateOfBirth)
		if (files && files[0] && isValidFileType(files[0].type)) {
			console.log('files')
			formData.append('image', files[0])
		}
		try {
			const response = await fetch('http://localhost:80/register', {
				method: 'POST',
				body: formData,

				credentials: 'include',
			})
			if (!response.ok) {
				const errorMessage = await response.text()
				msg = errorMessage // Set the error message to the msg variable
				throw new Error(`${errorMessage}`)
			}
			createWebSocket()
		} catch (error) {
			throw error
		}
	}
	$: if (files && files[0]) console.log(isValidFileType(files[0].type))
	$: console.log(files)
</script>

<main>
	<div class=" h-3/4 w-4/12 select-none" in:fade|global={{ duration: 500 }} out:fly|global={{ duration: 300, x: 800 }}>
		<div
			class="h-10 bg-gradient-to-t from-blue-500 to-blue-700 flex justify-center items-center border-2 rounded border-blue-950"
		>
			<h2 class="text-white font-extrabold">Register to Windows</h2>
		</div>
		<div class="h-30 bg-gray-500 logo flex justify-center border-b-2 border-stone-700">
			<img src={xplogo} class="h-20" alt="XP LOGO" />
		</div>
		<div
			class="h-100 p-5 bg-gray-500 flex flex-col items-center justify-center gap-5 border-b-2 border-gray-500 rounded-b-lg"
		>
			<div class="flex items-center gap-2">
				<label for="FirstName" class="w-24 text-right">First Name</label>
				<input
					bind:value={FirstName}
					type="first_name"
					placeholder="Type here"
					class="input w-full max-w-xs h-8 bg-white focus:outline-none"
				/>
				<label for="LastName" class="w-24 text-right">Last Name</label>
				<input
					bind:value={LastName}
					type="last_name"
					placeholder="Type here"
					class="input w-full max-w-xs h-8 bg-white focus:outline-none"
				/>
			</div>
			<div class="flex items-center gap-2">
				<label for="Email" class="w-24 text-right">Email</label>
				<input
					bind:value={Email}
					type="email"
					placeholder="Type here"
					class="input w-full max-w-xs h-8 bg-white focus:outline-none"
				/>
			</div>
			<div class="flex items-center gap-2">
				<label for="Password" class="w-24 text-right">Password</label>
				<input
					bind:value={Password}
					type="password"
					placeholder="Type here"
					class="input w-full max-w-xs h-8 bg-white focus:outline-none"
				/>
			</div>
			<div class="flex items-center gap-2">
				<label for="DateOfBirth" class="w-24 text-right">Date of Birth</label>
				<input
					bind:value={DateOfBirth}
					type="date"
					placeholder="Type here"
					class="input w-full max-w-xs h-8 bg-white focus:outline-none"
				/>
			</div>
			<div class="flex items-center gap-2">
				<label for="Avatar" class="w-24 text-right">Avatar {'(optional)'}</label>
				<input bind:files type="file" accept=".png, .jpg, .jpeg, .gif" class="w-full max-w-xs h-8 focus:outline-none" />
			</div>
			<div class="flex items-center gap-2">
				<label for="Nickname" class="w-24 text-right">Nickname {'(optional)'}</label>
				<input
					bind:value={Nickname}
					type="nickname"
					placeholder="Type here"
					class="input w-full max-w-xs h-8 focus:outline-none"
				/>
			</div>
			<div class="flex items-center gap-2">
				<label for="AboutMe" class="w-24 text-right">About {'(optional)'}</label>
				<input
					bind:value={AboutMe}
					type="about"
					placeholder="Type here"
					class="input w-full max-w-xs h-8 bg-white focus:outline-none"
				/>
			</div>
			<div class="flex flex-col items-center">
				<div class="flex flex-row">
					<button class="btn mr-3" on:click={register}>REGISTER</button>
					<button
						class="btn"
						on:click={() => {
							LoginPage.set('login')
						}}>TO LOGIN</button
					>
				</div>
				{#if registering}
					<h2 class="text-xl text-red-700 font-extrabold">{msg}</h2>
				{/if}
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
