<script>
	import { slide } from 'svelte/transition'
	import { currentUserGroups, eventStore, currentUser } from '../stores/user'
	import { socket } from '../ws'
	import { createEventDispatcher } from 'svelte'
	import { tr } from 'svelty-picker/i18n'
	import { formatTime } from '../utils'
	import MdCheck from 'svelte-icons/md/MdCheck.svelte'
	import MdCancel from 'svelte-icons/md/MdCancel.svelte'
	const dispatch = createEventDispatcher()

	export let group

	let owner = group.members[0]
	let members = group.members.slice(1)

	function joinGroup() {
		if (memberStatus) return

		const grp = group
		grp.status = 'requested'
		$currentUserGroups = [...$currentUserGroups, grp]
		found = null

		const data = {
			action: 'join_group',
			data: {
				id: group.id,
				creator_id: group.creator_id,
			},
		}
		socket.send(JSON.stringify(data))
	}

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
		const eventToAdd = group.events.find((item) => item.id == eventid)

		if (!$eventStore?.length) {
			$eventStore = []
		}
		$eventStore = [...$eventStore, eventToAdd]
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

	$: active = 'main'
	$: console.log(group)
	$: found = $currentUserGroups.find((grp) => grp.id == group.id)
	$: memberStatus = found ? found.status : null
</script>

<div class="border-4 border-red-400 mt-4 p-4 rounded">
	<h2 class="text-center text-4xl font-extrabold">{group.title}</h2>
	<p class="text-center">{group.description}</p>
	<p class="text-center">
		Created by <span
			class="font-extrabold text-primary hover:cursor-pointer"
			on:click={() => dispatch('user', owner.id)}>{owner.first_name} {owner.last_name}</span
		>
	</p>
	<div class="flex gap-4 mt-4 items-center" data-theme="dracula">
		<div class="flex flex-col w-full items-center">
			{#if memberStatus == 'joined'}
				<div class="w-full flex justify-center items-center">
					<div class="join grid grid-cols-2 bg-red-200 border-none">
						{#if group.events && group.events.length > 0}
							<button
								on:click={() => (active = 'main')}
								class="join-item btn btn-outline border-none text-base-200 hover:bg-red-100 {active == 'main'
									? 'bg-red-300'
									: ''}">Main page</button
							>

							<button
								on:click={() => (active = 'events')}
								class="join-item btn btn-outline border-none text-base-200 hover:bg-red-100 {active == 'events'
									? 'bg-red-300'
									: ''}">Events</button
							>
						{/if}
					</div>
				</div>
				<div class="flex gap-2 mb-4">
					<button class="btn w-40 mt-10" on:click={() => dispatch('post', group.id)}>View our posts</button>
					<button class="btn w-40 mt-10" on:click={() => dispatch('create_event', group.id)}>Create new event</button>
				</div>
				{#if active == 'main'}
					{#if group?.members}
						<h2 class="text-lg font-bold mt-10">Our members:</h2>
						{#each group.members as member}
							<h2 class="font-extrabold hover:cursor-pointer" on:click={() => dispatch('user', member.id)}>
								{member.first_name}
								{member.last_name}
							</h2>
						{/each}
					{/if}
				{:else if active == 'events'}
					<table class="table">
						<!-- head -->
						<thead>
							<tr class="text-black">
								<th />
								<th>Name</th>
								<th>Description</th>
								<th>Created By</th>
								<th>Start</th>
								<th>End</th>
								<th>Going</th>
							</tr>
						</thead>
						<tbody>
							{#each group.events as event, index (event.id)}
								<tr
									class="hover:bg-red-200 hover:cursor-pointer {eventOver(event.event_end) ? 'line-through' : ''}"
									on:click={() => dispatch('event', event)}
								>
									<th>{index + 1}</th>
									<td class="cursor-pointer font-extrabold">{event.title}</td>
									<td>{event.content}</td>
									<td>{event.user.first_name} {event.user.last_name}</td>

									<td>{formatTime(event.event_start)}</td>
									<td>{formatTime(event.event_end)}</td>
									<td
										>{event?.responses ? event.responses.filter((resp) => resp.response === 'going').length : 0}

										{#if !eventOver(event.event_end)}
											<div class="flex gap-2">
												{#if group.events.find((item) => item.id == event.id)?.responses && group.events
														.find((item) => item.id == event.id)
														.responses.find((item) => item.user_id == $currentUser.id && item.response == 'going')}
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
												{#if group.events.find((item) => item.id == event.id)?.responses && group.events
														.find((item) => item.id == event.id)
														.responses.find((item) => item.user_id == $currentUser.id && item.response == 'not going')}
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
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				{/if}
			{/if}
		</div>

		{#if memberStatus != 'joined'}
			<div class="flex w-full">
				{#if memberStatus == 'pending'}
					<button disabled class="btn btn-info disabled:text-black">Request sent, awaiting for answer.</button>
				{:else if !memberStatus}
					<div class="tooltip" data-tip="Send {owner.first_name} {owner.last_name} request to join {group.title}">
						<button class="disabled btn btn-info hover:btn-accent" on:click={joinGroup}>Request to join</button>
					</div>
				{/if}
				{#if memberStatus == 'requested'}
					<div class="alert w-fit text-info font-extrabold flex" transition:slide={{ duration: 300, axis: 'y' }}>
						<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-info shrink-0 w-6 h-6"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
							/></svg
						>
						<span
							>Request sent! {owner.first_name}
							{owner.last_name} has to accept it first. After joining a group you can see their events, posts and hang in
							their chatroom</span
						>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>
