<script>
	import { onMount } from 'svelte'
	import { slide } from 'svelte/transition'
	import { currentUserGroups } from '../stores/user'
	import SveltyPicker, { config } from 'svelty-picker'
	import { socket } from '../ws'
	config.weekStart = 1
	config.theme = 'Dark'

	let date
	let name
	let description
	export let selectedGroup
	$: console.log(selectedGroup)
	$: joinedGroups = $currentUserGroups ? $currentUserGroups.filter((group) => group.status == 'joined') : []

	$: disabled = date && date.length == 2 && description != '' && name != '' && selectedGroup > 0 ? false : true

	function createEvent() {
		const data = {
			action: 'new_event',
			data: {
				group_id: parseInt(selectedGroup),
				title: name,
				content: description,
				event_start: convertToUTC(date[0]),
				event_end: convertToUTC(date[1]),
			},
		}
		socket.send(JSON.stringify(data))
	}

	function convertToUTC(localDate) {
		const localDateObj = new Date(localDate)
		const timezoneOffset = localDateObj.getTimezoneOffset()
		const utcDate = new Date(localDateObj.getTime() + timezoneOffset)
		return utcDate.toISOString()
	}
</script>

<div
	class="bg-base-100 mt-4 rounded-3xl h-[100vh] p-10 flex gap-4 justify-center"
	in:slide|global={{ delay: 500, duration: 200, axis: 'y' }}
	out:slide|global={{ duration: 200, axis: 'x' }}
>
	<form class="flex flex-col" on:submit|preventDefault={createEvent}>
		<h2 class="text-primary text-4xl font-bold">Create event</h2>
		<select bind:value={selectedGroup} class="select select-bordered w-full max-w-xs" required>
			<option disabled selected>Select group</option>
			{#each joinedGroups as group}
				<option value={group.id}>{group.title}</option>
			{/each}
		</select>
		<div class="flex flex-col">
			<label for="name" class="text-white">Event name</label>
			<input required type="text" placeholder="Name" class="input input-bordered w-full max-w-xs" bind:value={name} />
		</div>
		<div class="flex flex-col">
			<label for="description" class="text-white">Event description</label>
			<input
				required
				type="text"
				placeholder="Description"
				class="input input-bordered w-full max-w-xs"
				bind:value={description}
			/>
		</div>
		<div class="flex flex-col">
			<label for="startDate" class="text-xl">Event duration [START-END]</label>

			<SveltyPicker
				required={true}
				isRange
				mode="datetime"
				pickerOnly
				bind:value={date}
				format="yyyy-mm-dd hh:ii"
				displayFormat="yyyy-mm-dd hh:ii"
				startDate={new Date()}
			/>
		</div>

		<button class="btn" {disabled}>Create event</button>
	</form>
</div>
