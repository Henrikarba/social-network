<script>
	import { slide } from 'svelte/transition'
	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()

	const iecrash = new URL('../assets/iecrash.JPG', import.meta.url).href
	let crash = false
	let groupName
	let groupDescription

	$: invalid =
		groupName &&
		groupName.trim() != '' &&
		groupName.length > 3 &&
		groupDescription &&
		groupDescription.trim() != '' &&
		groupDescription.length > 5
			? true
			: false

	async function createGroup() {
		if (!invalid) return
		crash = true
		const formData = new FormData()
		formData.append('title', groupName)
		formData.append('description', groupDescription)
		const result = await handleNewGroup(formData)
		groupName = ''
		groupDescription = ''
		console.log(result)
		setTimeout(() => {
			dispatch('group', result.group_id)
		}, 3000)
	}

	async function handleNewGroup(formData) {
		try {
			const response = await fetch(`http://localhost:80/new/group`, {
				method: 'POST',
				body: formData,
				credentials: 'include',
			})
			if (!response.ok) {
				const errorMessage = await response.text()
				throw new Error(`Request failed: ${errorMessage}`)
			}
			return await response.json()
		} catch (error) {
			throw error
		}
	}
</script>

{#if crash}
	<div class="flex justify-center items-center h-screen w-screen">
		<div class="absolute inset-0 flex items-center justify-center">
			<img src={iecrash} alt="" />
		</div>
	</div>
{/if}

<main
	data-theme="dracula"
	class="bg-base-100 mt-4 rounded-3xl h-[100vh] p-10 flex gap-4"
	in:slide|global={{ delay: 500, duration: 200, axis: 'y' }}
	out:slide|global={{ duration: 200, axis: 'x' }}
>
	<div class="w-full flex flex-col gap-4">
		<h2 class="text-center w-full text-4xl text-accent">Create new group</h2>
		<form on:submit|preventDefault={createGroup}>
			<div class="flex flex-col items-center">
				<label for="name" class="text-xl text-primary mb-2">Group name</label>
				<input
					bind:value={groupName}
					type="text"
					placeholder="name"
					class="input input-bordered w-full max-w-xs text-slate-200"
				/>
			</div>
			<div class="flex flex-col items-center">
				<label for="name" class="text-xl text-primary mb-2">Description</label>
				<textarea
					bind:value={groupDescription}
					class="textarea textarea-bordered w-full max-w-xs text-slate-200 h-40"
					placeholder="What is your group about?"
				/>
			</div>
			<div class="flex justify-center">
				<button class="btn btn-success w-fit" disabled={!invalid}>Submit</button>
			</div>
		</form>
	</div>
</main>
