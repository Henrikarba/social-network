<script>
	import { slide } from 'svelte/transition'
	import { currentUserGroups } from '../stores/user'
	import { socket } from '../ws'
	import { createEventDispatcher } from 'svelte'
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
	<div class="flex gap-4 mt-4" data-theme="dracula">
		<div class="flex flex-col">
			{#if memberStatus == 'joined'}
				<button class="btn" on:click={() => dispatch('post', group.id)}>View our posts</button>
				{#if group?.members}
					<h2 class="text-lg font-bold mt-10">Our members:</h2>
					{#each group.members as member}
						<h2 class="font-extrabold hover:cursor-pointer" on:click={() => dispatch('user', member.id)}>
							{member.first_name}
							{member.last_name}
						</h2>
					{/each}
				{/if}
			{/if}
		</div>

		{#if memberStatus != 'joined'}
			{#if memberStatus == 'pending'}
				<button disabled class="btn btn-info disabled:text-black">Request sent, awaiting for answer.</button>
			{:else if !memberStatus}
				<div class="tooltip" data-tip="Send {owner.first_name} {owner.last_name} request to join {group.title}">
					<button class="disabled btn btn-info hover:btn-accent" on:click={joinGroup}>Request to join</button>
				</div>
			{/if}
			{#if memberStatus == 'requested'}
				<div class="alert w-1/2 text-info font-extrabold" transition:slide={{ duration: 300, axis: 'y' }}>
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
						{owner.last_name} has to accept it first. After joining a group you can see their events, posts and hang in their
						chatroom</span
					>
				</div>
			{/if}
		{/if}
	</div>
</div>
