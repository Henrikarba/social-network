<script>
	import MdCheck from 'svelte-icons/md/MdCheck.svelte'
	import MdCancel from 'svelte-icons/md/MdCancel.svelte'
	import { socket } from '../ws'
	import { fly } from 'svelte/transition'
	import { flip } from 'svelte/animate'
	import { notificationsStore, currentUserFollowers } from '../stores/user'

	function handleAccept(notif) {
		$notificationsStore = $notificationsStore.filter((notiff) => notiff.id !== notif.id)

		let data = {
			action: 'accept',
			data: notif,
		}
		switch (notif.type) {
			case 'follow_request':
				$currentUserFollowers = $currentUserFollowers.map((follower) => {
					if (follower.id === notif.sender.id) {
						follower.status = 'accepted'
					}
					return follower
				})

				break
		}
		socket.send(JSON.stringify(data))
	}
	function handleReject(notif) {
		$notificationsStore = $notificationsStore.filter((notiff) => notiff.id !== notif.id)
		const data = {
			action: 'reject',
			data: notif,
		}
		socket.send(JSON.stringify(data))
	}

	$: console.log($notificationsStore)
	$: notifs = $notificationsStore
</script>

{#if notifs}
	<div data-theme="dracula" class="select-none" in:fly|global={{ duration: 500, delay: 5000, y: 9000 }}>
		{#each notifs as notif, index (notif.id)}
			<div
				animate:flip={{ duration: 500 }}
				out:fly={{ duration: 300, x: 800 }}
				class="p-4 h-32 w-96 bg-base-200 fixed mr-4 mt-4 rounded-3xl z-[9999999]"
				style="top:{index * 9}rem; right:0; "
			>
				<div class="flex items-center">
					{#if notif.sender.avatar && notif.sender.avatar != ''}
						<img class="w-8 h-8" src="http://localhost:80/images/{notif.sender.avatar}" alt="" />
					{/if}
					<p class="text-accent font-extrabold">
						{notif.sender.first_name}
						{notif.sender.last_name}
						<span class="text-info font-medium">{notif.message}</span>
						{#if notif?.group}
							<span class="text-primary">{notif.group.title}</span>
						{/if}
					</p>
				</div>
				{#if notif.type == 'join_request' || notif.type == 'follow_request' || notif.type == 'group_join_request'}
					<div class="flex justify-around mt-2">
						<div class="tooltip tooltip-success font-bold uppercase" data-tip="Accept">
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<!-- svelte-ignore a11y-no-static-element-interactions -->
							<div
								on:click={() => handleAccept(notif)}
								class="w-8 h-8 text-success hover:cursor-pointer hover:scale-125 transform
                                        transition duration-500"
							>
								<MdCheck />
							</div>
						</div>
						<div class="tooltip tooltip-error font-bold uppercase" data-tip="Reject">
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<!-- svelte-ignore a11y-no-static-element-interactions -->
							<div
								on:click={() => handleReject(notif)}
								class="w-8 h-8 text-error hover:cursor-pointer hover:scale-125 transform
                                        transition duration-500"
							>
								<MdCancel />
							</div>
						</div>
					</div>
				{:else}
					<div class="flex justify-around mt-2">
						<div class="tooltip tooltip-success font-bold uppercase ml-4" data-tip="Mark as read">
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<!-- svelte-ignore a11y-no-static-element-interactions -->
							<div
								on:click={() => handleAccept(notif)}
								class="w-8 h-8 text-success hover:cursor-pointer hover:scale-125 transform
                                        transition duration-500"
							>
								<MdCheck />
							</div>
						</div>
					</div>
				{/if}
			</div>
		{/each}
	</div>
{/if}
