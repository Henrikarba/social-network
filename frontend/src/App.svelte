<script>
  import { fade } from "svelte/transition";
  import { onMount } from "svelte";
  // WS
  import { isAuthenticated, socket, createWebSocket } from "./ws";
  import { messagesStore } from "./stores/chat";

  // Components
  import Footer from "./lib/Footer.svelte";
  import Shortcut from "./lib/Shortcut.svelte";
  import IE from "./lib/IE.svelte";
  import MSN from "./lib/MSN.svelte";
  import Login from "./lib/Login.svelte";
  import Register from "./lib/Register.svelte";
  import Milf from "./lib/Milf.svelte";
  import Notification from "./lib/Notification.svelte";
  import Chat from "./lib/Chat.svelte";
  import { LoginPage } from "./stores/signin";

  const msnUrl = new URL("./assets/msn.png", import.meta.url).href;
  const ieUrl = new URL("./assets/ie.png", import.meta.url).href;

  let ieOpen = false;
  let msnOpen = false;
  function openIE() {
    ieOpen = !ieOpen;
  }

  function openMSN() {
    msnOpen = !msnOpen;
  }

  let chatOpen = false;
  let chatType;
  let chatID;
  let groupname;
  function openChat(event) {
    chatType = event.detail.type;
    chatID = event.detail.id;
    if (chatType == "group") {
      groupname = event.detail.name;
    }
    if (
      $messagesStore &&
      chatType == "regular" &&
      $messagesStore.some((item) => item.sender_id == chatID)
    ) {
      $messagesStore = $messagesStore.filter(
        (item) => item.sender_id != chatID
      );
    }

    if (!chatOpen) chatOpen = true;
    focusElement(2);
  }

  $: elements = [
    { name: "msn", z: 50 },
    { name: "ie", z: 50 },
    { name: "chat", z: 50 },
  ];

  function focusElement(index) {
    elements[index].z =
      elements.reduce((maxZ, el) => Math.max(maxZ, el.z), 0) + 1;

    elements.forEach((el, i) => {
      if (i !== index) {
        el.z = Math.max(0, el.z - 1);
      }
    });
  }

  let loading = true;
  let authenticated = false;

  onMount(() => {
    createWebSocket();
    setTimeout(() => {
      loading = false;
    }, 3000);
  });

  function bsod() {
    setTimeout(() => {
      pizdec = false;
      loading = true;
    }, 1500);
    setTimeout(() => {
      loading = false;
    }, 3000);
  }

  let pizdec = false;

  $: authenticated = $isAuthenticated;
  $: if (pizdec) bsod();
</script>

<Notification />

{#if chatOpen}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div on:click={() => focusElement(2)}>
    <Chat
      {groupname}
      type={chatType}
      id={chatID}
      z={elements[2].z}
      on:last={() => focusElement(2)}
      on:close={() => (chatOpen = !chatOpen)}
    />
  </div>
{/if}

<main style={pizdec ? "transform: rotate(180deg)" : ""}>
  {#if authenticated && !loading}
    <Milf on:rotate={() => (pizdec = true)} />
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div on:click={() => focusElement(0)}>
      <Shortcut imgurl={msnUrl} left={300} on:open={openMSN}>MSN</Shortcut>
    </div>

    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div on:click={() => focusElement(1)}>
      <Shortcut imgurl={ieUrl} left={200} on:open={openIE}
        >Internet Explorer</Shortcut
      >
    </div>
    <Footer on:bsod={bsod} />

    {#if ieOpen}
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <div on:click={() => focusElement(1)}>
        <IE
          {ieUrl}
          on:close={openIE}
          on:last={() => focusElement(1)}
          z={elements[1].z}
        />
      </div>
    {/if}
    {#if msnOpen}
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <div on:click={() => focusElement(0)}>
        <MSN
          {msnUrl}
          on:close={openMSN}
          on:last={() => focusElement(0)}
          on:chat={openChat}
          z={elements[0].z}
        />
      </div>
    {/if}
  {/if}

  {#if !loading && !authenticated}
    {#if $LoginPage === "login"}
      <Login />
    {:else}
      <Register />
    {/if}
  {/if}
  {#if loading}
    <div out:fade={{ duration: 300 }} class="loader" />
  {/if}
</main>

<style>
  .loader {
    background-image: url("./assets/loading.gif");
    position: absolute;
    width: 100%;
    height: 100vh;
    z-index: 99999999;
    background-position: center center;
    background-size: 100%;
    background-attachment: fixed;
  }

  :global(.z-top) {
    z-index: 333;
  }

  :global(.z-low) {
    z-index: 125;
  }
</style>
