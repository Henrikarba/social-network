<script>
  import { currentUser } from "../stores/user";
  import { socket } from "../ws";
  // Svelte
  import { createEventDispatcher } from "svelte";
  const dispatch = createEventDispatcher();
  import { currentUserFollowing, currentUserGroups } from "../stores/user";
  import { postsStore } from "../stores/post";
  import { formatDateTime } from "../utils";
  export let profile;

  function makeFollowRequest(id) {
    if (!$currentUserFollowing) {
      $currentUserFollowing = [];
    }

    if (!followers) {
      followers = [];
    }

    if (!privateProfile) {
      $currentUserFollowing = [
        ...$currentUserFollowing,
        { ...profile.user, status: "accepted" },
      ];
      followers = [...followers, $currentUser];
    } else if (privateProfile) {
      $currentUserFollowing = [
        ...$currentUserFollowing,
        { ...profile.user, status: "pending" },
      ];
    }
    const data = {
      action: "follow_request",
      data: {
        id: parseInt(id),
      },
    };
    console.log("SENDING....");
    socket.send(JSON.stringify(data));
  }

  // Determine if the current user is following the profile
  let isFollowing = "";
  $: if ($currentUserFollowing) {
    const follower = $currentUserFollowing.find(
      (follower) => follower.id === profile.user.id
    );

    if (follower) {
      isFollowing = follower.status === "accepted" ? "accepted" : "pending";
    }
  }
  $: privateProfile = !profile.user?.email;
  let followers;
  let following;
  if (profile?.followers) {
    followers = profile.followers;
  }
  if (profile?.following) {
    following = profile.following;
  }
  let inviteableGroups;
  $: joinedGroups = $currentUserGroups
    ? $currentUserGroups.filter((group) => group.status == "joined")
    : [];
  $: if (joinedGroups && joinedGroups.length > 0) {
    inviteableGroups = joinedGroups.filter((item) => {
      return !item.member_ids.some((elem) => elem == profile.user.id);
    });
  }
  $: console.log(inviteableGroups);
  let selectedGroup;
  $: invites = [];
  function inviteToGroup() {
    const data = {
      action: "group_join_invite",
      data: {
        user_id: profile.user.id,
        id: parseInt(selectedGroup),
      },
    };
    $currentUserGroups = $currentUserGroups.map((item) => {
      if (item.id == selectedGroup) {
        item.member_ids = [...item.member_ids, profile.user.id];
      }
      return item;
    });
    socket.send(JSON.stringify(data));
    invites = [
      ...invites,
      `Invited ${profile.user.first_name} ${profile.user.last_name} to ${
        inviteableGroups.find((item) => item.id == selectedGroup).title
      }`,
    ];
  }
  $: console.log(selectedGroup);
  $: userPosts = $postsStore.filter((item) => item.user_id == profile.user.id);
</script>

<div class="flex items-center mt-10 border-4 flex-col p-6">
  <h2>Viewing the profile of:</h2>
  {#if profile.user?.avatar}
    <img
      class="w-16"
      src="http://localhost:80/images/{profile.user.avatar}"
      alt=""
    />
  {/if}
  <h2 class="text-4xl font-extrabold">
    {profile.user.first_name}
    {profile.user.last_name}
  </h2>

  {#if isFollowing == "" || (isFollowing == "pending" && privateProfile)}
    <div>
      <p>
        This person is using technologies and such to stay private, so no more
        info for you.
      </p>
      <br />
      <p>
        You may try following him or her and when they accept, you might see
        something more interesting
      </p>
      <br />
      <p class="font-mono">no promises though</p>
    </div>
  {:else}
    <h2>{profile.user.email}</h2>
    {#if profile.user?.about_me}
      <h2>{profile.user.about_me}</h2>
    {/if}
    <h2>Born on {formatDateTime(profile.user.date_of_birth)}</h2>
  {/if}
  {#if userPosts && userPosts.length > 0}
    <h2 class="font-bold text-xl">Posts created by me:</h2>
    {#each userPosts as post}
      <h2
        class="text-primary hover:cursor-pointer"
        on:click={() => dispatch("singlePost", post.post_id)}
      >
        {post.title}
      </h2>
    {/each}
  {/if}
  {#if followers && followers.length > 0}
    <h2 class="mt-4 border-t-4">My followers:</h2>
    <div class="border-2 border-red-950 p-2 flex gap-2">
      {#each followers as follower, index (follower.id)}
        <h2
          class="text-orange-500 font-extrabold cursor-pointer"
          on:click={() => dispatch("user", follower.id)}
        >
          {follower.first_name}
          {follower.last_name}
        </h2>
      {/each}
    </div>
  {/if}
  {#if following && following.length > 0}
    <h2 class="mt-4 border-t-4">I follow:</h2>
    <div class="border-2 border-red-950 p-2 flex gap-2">
      {#each following as follower, index (follower.id)}
        <h2
          class="text-orange-500 font-extrabold cursor-pointer"
          on:click={() => dispatch("user", follower.id)}
        >
          {follower.first_name}
          {follower.last_name}
        </h2>
      {/each}
    </div>
  {/if}
  {#if isFollowing == ""}
    <button on:click={() => makeFollowRequest(profile.user.id)} class="btn"
      >{privateProfile ? "Send follow request" : "Start following"}</button
    >
  {:else if isFollowing == "pending"}
    <h2 class="p-6 border-2 rounded border-red-500">
      Follow request sent, waiting to hear back..
    </h2>
  {/if}
  {#if inviteableGroups && inviteableGroups.length > 0}
    <div class="mt-10">
      <h2>Invite user to your group(s)</h2>
      <select name="group" bind:value={selectedGroup}>
        <option disabled selected value="0">Select</option>
        {#each inviteableGroups as group, index (group.id)}
          <option value={group.id}>{group.title}</option>
        {/each}
      </select>
      <button
        class="btn"
        on:click={inviteToGroup}
        disabled={selectedGroup == 0 ? true : false}>Invite user</button
      >
    </div>
  {/if}
  {#if invites && invites.length > 0}
    {#each invites as invite}
      <h2 class="text-primary">{invite}!</h2>
    {/each}
  {/if}
</div>
