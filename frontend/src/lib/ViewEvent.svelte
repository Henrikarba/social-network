<script>
	import { formatTime } from '../utils'
	import MdCancel from 'svelte-icons/md/MdCancel.svelte'
	import MdCheck from 'svelte-icons/md/MdCheck.svelte'

	export let event
	import { currentUser } from '../stores/user'
	import { eventStore } from '../stores/user'
	import { socket } from '../ws'

	function eventOver(dateString) {
		const current = new Date()
		const eventDate = new Date(dateString)

		return current > eventDate
	}

	function handleEventResponse(eventid, response) {
		const data = {
			action: 'event_response',
			data: {
				event_id: eventid,
				response: response,
				user_id: $currentUser.id,
			},
		}
		socket.send(JSON.stringify(data))

		if (!$eventStore?.length) {
			$eventStore = []
		}
		$eventStore = [...$eventStore, event]
		group.events = group.events.map((item) => {
			if (item.id == eventid) {
				if (!item.responses) {
					item.responses = []
				}
				const userResponse = item.responses.find((resp) => resp.user_id == $currentUser.id)
				if (!userResponse) {
					item.responses.push(data.data)
				} else {
					userResponse.response = response
				}
			}
			return item
		})
	}
	$: console.log(event)
</script>

<div class="border-4 border-red-400 mt-4 p-4 rounded">
	<h2 class="text-center text-4xl">{event.title}</h2>
	<p class="text-center">by {event.user.first_name} {event.user.last_name}</p>

	<h2 class="text-xl">About this event:</h2>
	<h2><span class="font-bold">Description:</span> {event.content}</h2>
	{#if eventOver(event.event_end)}
		<h2>This event ended on {formatTime(event.event_end)}</h2>
	{:else}
		<h2>Event start: <span class="text-primary font-extrabold">{formatTime(event.event_start)}</span></h2>
		<h2>Event end: <span class="text-primary font-extrabold">{formatTime(event.event_end)}</span></h2>
	{/if}
	{#if !eventOver(event.event_end)}
		<div class="flex gap-2">
			{#if event.responses.find((item) => item.user_id == $currentUser.id && item.response == 'going')}
				<h2 class="text-primary">Going!</h2>
			{:else}
				<div class="tooltip" data-tip="Going">
					<button
						class="btn btn-sm btn-circle btn-success"
						on:click|preventDefault|stopPropagation={() => handleEventResponse(event.id, 'going')}
					>
						<MdCheck />
					</button>
				</div>
			{/if}
			{#if event.responses.find((item) => item.user_id == $currentUser.id && item.response == 'not going')}
				<h2 class="text-primary">Not going</h2>
			{:else}
				<div class="tooltip" data-tip="Not Going">
					<button
						class="btn btn-sm btn-circle btn-error"
						on:click|preventDefault|stopPropagation={() => handleEventResponse(event.id, 'not going')}
						><MdCancel /></button
					>
				</div>
			{/if}
		</div>
	{/if}
	{#if event.responses && event.responses.length > 0}
		{#if event.responses.filter((item) => item.response == 'going').length > 0}
			<h2>Marked as going:</h2>

			{#each event.responses.filter((item) => item.response == 'going') as response}
				<h2>{response.user.first_name} {response.user.last_name}</h2>
			{/each}
		{/if}
		{#if event.responses.filter((item) => item.response == 'not going').length > 0}
			<h2>Marked as not going:</h2>
			{#each event.responses.filter((item) => item.response == 'not going') as response}
				<h2>{response.user.first_name} {response.user.last_name}</h2>
			{/each}
		{/if}
	{/if}
</div>
