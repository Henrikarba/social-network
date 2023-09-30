<script>
	import { formatTime } from '../utils'
	export let event

	function eventOver(dateString) {
		const current = new Date()
		const eventDate = new Date(dateString)

		return current > eventDate
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
