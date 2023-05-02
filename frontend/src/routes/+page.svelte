<script lang="ts">
	// TODO: Finish style for the pages
	// TODO: Implement session creation

	import { onMount } from 'svelte';

	import { writable } from 'svelte/store';

	type Sessions = {
		id: number;
		sessionName: string;
	};

	type ValidateSessionSelected = {
		id: number;
		selected: boolean;
	};

	type Message = {
		where?: string;
		from: string | number;
		message: string;
	};

	const user = writable<number | null>(0 || null);
	const sessions = writable<Sessions[]>([]);

	onMount(async () => {
		if (sessionStorage.getItem('token') === null) {
			window.location.href = '/login';
		}
		await fetch('https://localhost:8080/api/getSessions', {
			method: 'GET',
			headers: {
				Authorization: sessionStorage.getItem('authToken') || ''
			}
		})
			.then((res) => res.json())
			.then((data) => {
				sessions.set(data.sessions);
				user.set(sessionStorage.getItem('token') ? Number(sessionStorage.getItem('token')) : null);
			})
			.catch((err) => console.log(err));
	});

	let message: Message;
	let id = writable<number>(0);
	let messages = writable<Message[]>([]);
	let socket: WebSocket;

	$: id.subscribe((val) => {
		if (val !== 0) {
			openConnection();
		}
	});

	const statusJoinModal = writable<ValidateSessionSelected>({} as ValidateSessionSelected);
	const sessionToken = writable<string>('');

	function openConnection(): void {
		socket = new WebSocket('wss://localhost:8080/wss/login', [$id.toString()]);
		socket.addEventListener('open', () => {
			socket.send(
				JSON.stringify({
					sessionToken: $sessionToken,
					id: $id
				})
			);
		});
		socket.addEventListener('message', function (event) {
			if (!event.data.includes('User authenticated')) {
				const msg = JSON.parse(event.data);
				if (msg.message === 'close') {
					socket.close();
					return;
				}
				message = { from: msg.from, message: msg.message } as Message;
				messages.update((msgs) => [...msgs, message]);
			}
		});
	}

	function sendMessage(): void {
		const sendbox = document.getElementById('sendbox') as HTMLInputElement;
		socket.send(
			JSON.stringify({
				from: $user,
				message: sendbox?.value as string
			})
		);
		sendbox!.value = '';
	}

	function interactOnMenu(sessionId: number) {
		if ($statusJoinModal.selected == true) {
			statusJoinModal.set({
				id: sessionId,
				selected: false
			});
		} else {
			statusJoinModal.set({
				id: sessionId,
				selected: true
			});
		}
	}
</script>

<div class="flex flex-col w-screen h-screen bg-zinc-950 p-1 text-white">
	<nav class="flex w-full h-12 justify-between pl-4 pr-4">
		<h1>Header</h1>
		<button
			on:click={() => {
				sessionStorage.removeItem('token');
				if (socket) socket.close();
				window.location.href = '/login';
			}}
		>
			Leave
		</button>
	</nav>
	<div class="flex w-full items-center justify-around h-full">
		<div class="flex flex-col w-1/6 h-full p-4">
			<h1>Sessions</h1>
			{#each $sessions as session}
				<div class="flex p-1">
					<h1><button on:click={() => interactOnMenu(session.id)}>{session.id} -</button></h1>
					<h1>
						<button on:click={() => interactOnMenu(session.id)}>{session.sessionName}</button>
					</h1>
					{#if session.id == $id}
						<button
							class="ml-3"
							on:click={() => {
								socket.send(
									JSON.stringify({
										from: session.id,
										message: 'close'
									})
								);
								id.set(0);
								messages.set([]);
							}}>leave</button
						>
					{/if}
				</div>
				{#if $statusJoinModal.selected == true && $statusJoinModal.id == session.id}
					<div class="flex flex-col">
						<input type="text" name="sessionToken" id="sessionToken" />
						<button
							on:click={() => {
								const sT = document.getElementById('sessionToken')?.value;
								sessionToken.set(sT);
								id.set(session.id);
							}}>Join</button
						>
					</div>
				{/if}
			{/each}
			<!-- <div class="flex mt-4"> -->
			<!-- 	<h1>Add</h1> -->
			<!-- 	<h1>Create</h1> -->
			<!-- </div> -->
		</div>
		<div class="flex flex-col bg-zinc-900 w-5/6 h-full justify-around p-4 relative">
			<div class="flex flex-row p-4 w-full h-full pr-6">
				{#if $id == 0}
					<h1 class="text-2xl justify-center align-middle">Select a session to join</h1>
				{:else}
					<div class="flex flex-col w-1/2 h-full items-start">
						{#each $messages as msg}
							{#if msg.from == $user}
								<div class="flex justify-end rounded-md p-2 h-24" />
								<!-- <br /> -->
							{/if}
							{#if msg.from != $user}
								<div
									class="flex items-center bg-indigo-800 min-w-[20%] max-w-prose rounded-md p-2 h-24 max-h-full"
								>
									<p class="text-white break-words">{msg.message}</p>
								</div>
								<br />
							{/if}
							{#if msg.from == $user}
								<div class="flex justify-end rounded-md p-2 h-24" />
								<!-- <br /> -->
							{/if}
							<!-- <div class="flex justify-end rounded-md p-2 h-24" /> -->
						{/each}
					</div>
					<div />
					<div class="flex flex-col w-1/2 h-full items-end min-w-1/4">
						{#each $messages as msg}
							{#if msg.from != $user}
								<div class="flex justify-end rounded-md p-2 h-24" />
								<!-- <br /> -->
							{/if}
							{#if msg.from == $user}
								<div
									class="flex items-center bg-indigo-800 min-w-[20%] max-w-prose rounded-md p-2 h-24 max-h-full"
								>
									<p class="text-white break-words">{msg.message}</p>
								</div>
								<br />
							{/if}
							{#if msg.from != $user}
								<div class="flex justify-end rounded-md p-2 h-24" />
								<!-- <br /> -->
							{/if}
						{/each}
					</div>
				{/if}
			</div>
			<div class="flex w-full p-2 bg-cyan-900">
				<textarea
					name="sendbox"
					id="sendbox"
					value=""
					class="flex w-full h-14 rounded-md text-black"
				/>
				<button on:click={sendMessage}>Send</button>
			</div>
		</div>
	</div>
</div>
